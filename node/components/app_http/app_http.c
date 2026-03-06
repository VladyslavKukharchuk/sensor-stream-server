#include "app_http.h"
#include "esp_http_client.h"
#include "esp_crt_bundle.h"
#include "esp_log.h"
#include "esp_mac.h"
#include "cJSON.h"
#include <string.h>
#include <time.h>

static const char *TAG = "app_http";

typedef struct {
    char buffer[512];
    int len;
} http_response_data_t;

static esp_err_t _http_event_handler(esp_http_client_event_t *evt) {
    if (evt->event_id == HTTP_EVENT_ON_DATA) {
        http_response_data_t *res_data = (http_response_data_t *)evt->user_data;
        if (res_data && (res_data->len + evt->data_len < sizeof(res_data->buffer))) {
            memcpy(res_data->buffer + res_data->len, evt->data, evt->data_len);
            res_data->len += evt->data_len;
            res_data->buffer[res_data->len] = '\0';
        }
    }
    return ESP_OK;
}

bool register_device(const char* server_url, char* out_device_id, size_t size) {
    uint8_t mac[6];
    esp_read_mac(mac, ESP_MAC_WIFI_STA);
    char mac_str[18];
    snprintf(mac_str, sizeof(mac_str), "%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);

    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "mac", mac_str);
    char *post_data = cJSON_PrintUnformatted(root);

    char url[256];
    snprintf(url, sizeof(url), "%s/api/v1/devices", server_url);

    http_response_data_t response_data = { .len = 0 };
    esp_http_client_config_t config = {
        .url = url,
        .method = HTTP_METHOD_POST,
        .crt_bundle_attach = esp_crt_bundle_attach,
        .skip_cert_common_name_check = true,
        .event_handler = _http_event_handler,
        .user_data = &response_data,
    };
    
    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_post_field(client, post_data, strlen(post_data));

    bool success = false;
    if (esp_http_client_perform(client) == ESP_OK) {
        int status = esp_http_client_get_status_code(client);
        if (status >= 200 && status < 300 && response_data.len > 0) {
            cJSON *resp = cJSON_Parse(response_data.buffer);
            cJSON *id_item = cJSON_GetObjectItem(resp, "id");
            if (cJSON_IsString(id_item)) {
                strncpy(out_device_id, id_item->valuestring, size);
                success = true;
            }
            cJSON_Delete(resp);
        }
    }
    esp_http_client_cleanup(client);
    cJSON_Delete(root);
    free(post_data);
    return success;
}

void send_measurement(const char* server_url, const char* device_id, float temp, float hum) {
    time_t now;
    struct tm timeinfo;
    time(&now);
    localtime_r(&now, &timeinfo);
    char timestamp[32];
    strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);

    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "device_id", device_id);
    cJSON_AddNumberToObject(root, "temperature", (int)(temp * 10 + 0.5) / 10.0);
    cJSON_AddNumberToObject(root, "humidity", (int)(hum * 10 + 0.5) / 10.0);
    cJSON_AddStringToObject(root, "timestamp", timestamp);
    char *post_data = cJSON_PrintUnformatted(root);

    char url[256];
    snprintf(url, sizeof(url), "%s/api/v1/measurements", server_url);

    esp_http_client_config_t config = {
        .url = url,
        .method = HTTP_METHOD_POST,
        .crt_bundle_attach = esp_crt_bundle_attach,
        .skip_cert_common_name_check = true,
    };
    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_post_field(client, post_data, strlen(post_data));

    esp_err_t err = esp_http_client_perform(client);
    if (err == ESP_OK) {
        ESP_LOGI(TAG, "POST Success. Status: %d", esp_http_client_get_status_code(client));
    }
    esp_http_client_cleanup(client);
    cJSON_Delete(root);
    free(post_data);
}
