#ifndef APP_WIFI_H
#define APP_WIFI_H

#include "esp_err.h"
#include "freertos/FreeRTOS.h"
#include "freertos/event_groups.h"

#define WIFI_CONNECTED_BIT BIT0

extern EventGroupHandle_t s_wifi_event_group;

void wifi_init_sta(const char* ssid, const char* password);

#endif
