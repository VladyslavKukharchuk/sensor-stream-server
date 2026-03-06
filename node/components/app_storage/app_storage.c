#include "app_storage.h"
#include "esp_spiffs.h"
#include "esp_log.h"
#include <stdio.h>
#include <string.h>

static const char *TAG = "app_storage";

void init_spiffs(void) {
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

bool read_device_id(char* device_id, size_t size) {
    FILE* f = fopen("/spiffs/device_id.txt", "r");
    if (f == NULL) return false;
    fgets(device_id, size, f);
    fclose(f);
    device_id[strcspn(device_id, "
")] = 0;
    return strlen(device_id) > 0;
}

void save_device_id(const char* device_id) {
    FILE* f = fopen("/spiffs/device_id.txt", "w");
    if (f != NULL) {
        fprintf(f, "%s", device_id);
        fclose(f);
    }
}
