#ifndef APP_STORAGE_H
#define APP_STORAGE_H

#include <stdbool.h>
#include <stddef.h>

void init_spiffs(void);
bool read_device_id(char* device_id, size_t size);
void save_device_id(const char* device_id);

#endif
