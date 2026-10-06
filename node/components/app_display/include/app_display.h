#ifndef APP_DISPLAY_H
#define APP_DISPLAY_H

#include "esp_err.h"

esp_err_t app_display_init(void);
esp_err_t app_display_show_startup(void);
esp_err_t app_display_show_reading(float temperature, float humidity);
esp_err_t app_display_show_sensor_error(void);

#endif
