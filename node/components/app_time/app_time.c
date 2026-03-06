#include "app_time.h"
#include "esp_netif_sntp.h"
#include "esp_log.h"
#include <time.h>
#include <stdlib.h>

static const char *TAG = "app_time";

void setup_time(void) {
    ESP_LOGI(TAG, "Initializing SNTP");
    esp_sntp_config_t config = ESP_NETIF_SNTP_DEFAULT_CONFIG("pool.ntp.org");
    esp_netif_sntp_init(&config);
    
    // Set timezone to Ukraine (EET-2EEST)
    setenv("TZ", "EET-2EEST,M3.5.0/3,M10.5.0/4", 1);
    tzset();
}
