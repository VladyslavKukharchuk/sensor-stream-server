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
#include "esp_crt_bundle.h"
#include "esp_mac.h"
#include "cJSON.h"
#include "dht.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/event_groups.h"

#include "secrets.h"

static const char *TAG = "SENSOR_NODE";

#define DHT_GPIO 4
#define DHT_TYPE DHT_TYPE_AM2301

#define WIFI_CONNECTED_BIT BIT0
static EventGroupHandle_t s_wifi_event_group;
static char device_id[64] = {0};

// --- WiFi ---
static void event_handler(void* arg, esp_event_base_t event_base, int32_t event_id, void* event_data) {
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {
        esp_wifi_connect();
    } else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED) {
        wifi_event_sta_disconnected_t* event = (wifi_event_sta_disconnected_t*) event_data;
        ESP_LOGW(TAG, "WiFi Disconnected. Reason: %d", event->reason);
        esp_wifi_connect();
    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {
        ip_event_got_ip_t* event = (ip_event_got_ip_t*) event_data;
        ESP_LOGI(TAG, "got ip:" IPSTR, IP2STR(&event->ip_info.ip));
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
      .partition_label = "storage",
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

// --- Registration ---
bool register_device() {
    uint8_t mac[6];
    esp_read_mac(mac, ESP_MAC_WIFI_STA);
    char mac_str[18];
    snprintf(mac_str, sizeof(mac_str), "%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);

    ESP_LOGI(TAG, "Registering device with MAC: %s", mac_str);

    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "mac", mac_str);
    char *post_data = cJSON_PrintUnformatted(root);

    esp_http_client_config_t config = {
        .url = SERVER_URL "/api/v1/devices",
        .method = HTTP_METHOD_POST,
        .crt_bundle_attach = esp_crt_bundle_attach,
        .skip_cert_common_name_check = true,
    };
    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_post_field(client, post_data, strlen(post_data));

    bool success = false;
    esp_err_t err = esp_http_client_perform(client);
    if (err == ESP_OK) {
        int status = esp_http_client_get_status_code(client);
        if (status >= 200 && status < 300) {
            char response_buffer[256];
            int len = esp_http_client_read_response(client, response_buffer, sizeof(response_buffer) - 1);
            if (len > 0) {
                response_buffer[len] = '\0';
                cJSON *resp = cJSON_Parse(response_buffer);
                cJSON *id_item = cJSON_GetObjectItem(resp, "id");
                if (cJSON_IsString(id_item)) {
                    strcpy(device_id, id_item->valuestring);
                    save_device_id(device_id);
                    ESP_LOGI(TAG, "Registered successfully. Assigned ID: %s", device_id);
                    success = true;
                }
                cJSON_Delete(resp);
            }
        } else {
            ESP_LOGE(TAG, "Registration failed with status: %d", status);
        }
    } else {
        ESP_LOGE(TAG, "Registration request failed: %s", esp_err_to_name(err));
    }

    esp_http_client_cleanup(client);
    cJSON_Delete(root);
    free(post_data);
    return success;
}

// --- HTTP Measurements ---
void send_measurement(float temp, float hum) {
    char timestamp[32];
    get_timestamp(timestamp, sizeof(timestamp));

    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "device_id", device_id);
    cJSON_AddNumberToObject(root, "temperature", (int)(temp * 10) / 10.0);
    cJSON_AddNumberToObject(root, "humidity", (int)(hum * 10) / 10.0);
    cJSON_AddStringToObject(root, "timestamp", timestamp);
    char *post_data = cJSON_PrintUnformatted(root);

    esp_http_client_config_t config = {
        .url = SERVER_URL "/api/v1/measurements",
        .method = HTTP_METHOD_POST,
        .crt_bundle_attach = esp_crt_bundle_attach,
        .skip_cert_common_name_check = true,
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

// --- Sensor Task ---
void sensor_task(void *pvParameters) {
    float temp, hum;

    while(1) {
        esp_err_t res = dht_read_float_data(DHT_TYPE, DHT_GPIO, &hum, &temp);
        if (res == ESP_OK) {
            ESP_LOGI(TAG, "Reading: T=%.1f C, H=%.1f%%", temp, hum);
            send_measurement(temp, hum);
        } else {
            ESP_LOGE(TAG, "DHT Sensor Error: %s (%d)", esp_err_to_name(res), res);
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
        ESP_LOGI(TAG, "Device ID not found in flash. Registering...");
        if (!register_device()) {
            ESP_LOGE(TAG, "Could not register device. Using fallback ID.");
            strcpy(device_id, "fallback-id");
        }
    }
    ESP_LOGI(TAG, "Device ID: %s", device_id);

    xTaskCreate(sensor_task, "sensor_task", 4096, NULL, 5, NULL);
}
