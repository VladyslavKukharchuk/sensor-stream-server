#include <stdio.h>
#include <string.h>
#include "esp_system.h"
#include "esp_wifi.h"
#include "esp_event.h"
#include "esp_log.h"
#include "nvs_flash.h"
#include "esp_http_client.h"
#include "esp_netif_sntp.h"
#include "esp_spiffs.h"
#include "cJSON.h"
#include "dht.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/event_groups.h"

#include "secrets.h"

static const char *TAG = "SENSOR_NODE";

#define DHT_GPIO 4
#define DHT_TYPE DHT_TYPE_DHT22

#define WIFI_CONNECTED_BIT BIT0
static EventGroupHandle_t s_wifi_event_group;
static char device_id[64] = {0};

// --- WiFi ---
static void event_handler(void* arg, esp_event_base_t event_base, int32_t event_id, void* event_data) {
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {
        esp_wifi_connect();
    } else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED) {
        esp_wifi_connect();
        ESP_LOGI(TAG, "retry to connect to the AP");
    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {
        xEventGroupSetBits(s_wifi_event_group, WIFI_CONNECTED_BIT);
    }
}

void wifi_init_sta(void) {
    s_wifi_event_group = xEventGroupCreate();
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    esp_netif_create_default_wifi_sta();
    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));

    esp_event_handler_instance_t instance_any_id;
    esp_event_handler_instance_t instance_got_ip;
    ESP_ERROR_CHECK(esp_event_handler_instance_register(WIFI_EVENT, ESP_EVENT_ANY_ID, &event_handler, NULL, &instance_any_id));
    ESP_ERROR_CHECK(esp_event_handler_instance_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &event_handler, NULL, &instance_got_ip));

    wifi_config_t wifi_config = {
        .sta = {
            .ssid = WIFI_SSID,
            .password = WIFI_PASSWORD,
        },
    };
    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());
    ESP_LOGI(TAG, "wifi_init_sta finished.");
}

// --- Time ---
void setup_time() {
    esp_sntp_config_t config = ESP_NETIF_SNTP_DEFAULT_CONFIG("pool.ntp.org");
    esp_netif_sntp_init(&config);
    setenv("TZ", "EET-2EEST,M3.5.0/3,M10.5.0/4", 1);
    tzset();
}

void get_timestamp(char *buf, size_t len) {
    time_t now;
    struct tm timeinfo;
    time(&now);
    localtime_r(&now, &timeinfo);
    strftime(buf, len, "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
}

// --- Storage ---
void init_spiffs() {
    esp_vfs_spiffs_conf_t conf = {
      .base_path = "/spiffs",
      .partition_label = NULL,
      .max_files = 5,
      .format_if_mount_failed = true
    };
    esp_err_t ret = esp_vfs_spiffs_register(&conf);
    if (ret != ESP_OK) {
        ESP_LOGE(TAG, "Failed to mount SPIFFS");
    }
}

bool read_device_id() {
    FILE* f = fopen("/spiffs/device_id.txt", "r");
    if (f == NULL) return false;
    fgets(device_id, sizeof(device_id), f);
    fclose(f);
    return strlen(device_id) > 0;
}

void save_device_id(const char* id) {
    FILE* f = fopen("/spiffs/device_id.txt", "w");
    if (f != NULL) {
        fprintf(f, "%s", id);
        fclose(f);
    }
}

// --- HTTP ---
void send_measurement(float temp, float hum) {
    char timestamp[32];
    get_timestamp(timestamp, sizeof(timestamp));

    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "device_id", device_id);
    cJSON_AddNumberToObject(root, "temperature", temp);
    cJSON_AddNumberToObject(root, "humidity", hum);
    cJSON_AddStringToObject(root, "timestamp", timestamp);
    char *post_data = cJSON_PrintUnformatted(root);

    esp_http_client_config_t config = {
        .url = SERVER_URL "/api/v1/measurements",
        .method = HTTP_METHOD_POST,
    };
    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_post_field(client, post_data, strlen(post_data));

    esp_err_t err = esp_http_client_perform(client);
    if (err == ESP_OK) {
        ESP_LOGI(TAG, "HTTP POST Status = %d", esp_http_client_get_status_code(client));
    } else {
        ESP_LOGE(TAG, "HTTP POST request failed: %s", esp_err_to_name(err));
    }
    esp_http_client_cleanup(client);
    cJSON_Delete(root);
    free(post_data);
}

// --- Main Task ---
void sensor_task(void *pvParameters) {
    float temp, hum;
    while(1) {
        esp_err_t res = dht_read_float_data(DHT_TYPE, DHT_GPIO, &hum, &temp);
        if (res == ESP_OK) {
            ESP_LOGI(TAG, "Reading: T=%.1f C, H=%.1f%%", temp, hum);
            send_measurement(temp, hum);
        } else {
            ESP_LOGE(TAG, "Could not read data from DHT sensor: %d", res);
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

    wifi_init_sta();
    xEventGroupWaitBits(s_wifi_event_group, WIFI_CONNECTED_BIT, pdFALSE, pdTRUE, portMAX_DELAY);
    
    setup_time();
    init_spiffs();
    
    if (!read_device_id()) {
        ESP_LOGI(TAG, "Device ID not found. Need registration logic here (similar to Arduino)");
        // Registration logic placeholder
        strcpy(device_id, "esp32-c6-idf-001"); 
    }

    xTaskCreate(sensor_task, "sensor_task", 4096, NULL, 5, NULL);
}
