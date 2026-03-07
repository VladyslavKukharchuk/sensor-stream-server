#include <stdio.h>
#include <string.h>
#include "esp_log.h"
#include "nvs_flash.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

#include "app_wifi.h"
#include "app_storage.h"
#include "app_http.h"
#include "app_time.h"
#include "dht.h"
#include "secrets.h"

static const char *TAG = "MAIN";

#define DHT_GPIO 4
#define DHT_TYPE DHT_TYPE_AM2301

static char device_id[64] = {0};

void sensor_task(void *pvParameters) {
    float temp, hum;

    while(1) {
        if (dht_read_float_data(DHT_TYPE, DHT_GPIO, &hum, &temp) == ESP_OK) {
            ESP_LOGI(TAG, "Reading: T=%.1f°C, H=%.1f%%", temp, hum);
            send_measurement(SERVER_URL, device_id, temp, hum);
        } else {
            ESP_LOGE(TAG, "Failed to read from DHT sensor");
        }
        vTaskDelay(pdMS_TO_TICKS(SEND_INTERVAL_SEC * 1000));
    }
}

void app_main(void) {
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND) {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);

    wifi_init_sta(WIFI_SSID, WIFI_PASSWORD);
    
    ESP_LOGI(TAG, "Waiting for WiFi...");
    xEventGroupWaitBits(s_wifi_event_group, WIFI_CONNECTED_BIT, pdFALSE, pdTRUE, portMAX_DELAY);

    setup_time();
    init_spiffs();

    if (!read_device_id(device_id, sizeof(device_id))) {
        ESP_LOGI(TAG, "Device not registered. Attempting registration...");
        if (!register_device(SERVER_URL, device_id, sizeof(device_id))) {
            ESP_LOGE(TAG, "CRITICAL ERROR: Registration failed. Device ID is mandatory.");
            ESP_LOGE(TAG, "System halted. Please check server connectivity.");

            return;
        }
    }
    
    ESP_LOGI(TAG, "Authenticated Device ID: %s", device_id);

    xTaskCreate(sensor_task, "sensor_task", 4096, NULL, 5, NULL);
}
