#ifndef APP_HTTP_H
#define APP_HTTP_H

#include <stdbool.h>
#include "esp_err.h"

bool register_device(const char* server_url, char* out_device_id, size_t size);
void send_measurement(const char* server_url, const char* device_id, float temp, float hum);

#endif
