#include "dht.h"
#include "sdkconfig.h"

#include <freertos/FreeRTOS.h>
#include <freertos/task.h>
#include <string.h>
#include <esp_log.h>
#include <ets_sys.h>
#include <esp_timer.h>
#include <driver/gpio.h>

#if HELPER_TARGET_IS_ESP32 || defined(CONFIG_IDF_TARGET_ESP32C6)
static portMUX_TYPE mux = portMUX_INITIALIZER_UNLOCKED;
#define PORT_ENTER_CRITICAL() portENTER_CRITICAL(&mux)
#define PORT_EXIT_CRITICAL() portEXIT_CRITICAL(&mux)
#else
#define PORT_ENTER_CRITICAL() 
#define PORT_EXIT_CRITICAL() 
#endif

static esp_err_t dht_await_pin_state(gpio_num_t pin, uint32_t timeout_us, int expected_pin_state)
{
    int64_t start = esp_timer_get_time();
    while (gpio_get_level(pin) != expected_pin_state)
    {
        if ((esp_timer_get_time() - start) > timeout_us) return ESP_ERR_TIMEOUT;
        ets_delay_us(1);
    }
    return ESP_OK;
}

static inline esp_err_t dht_fetch_data_direct(gpio_num_t pin, uint8_t data[5])
{
    // Phase B: Wait for sensor to pull low
    if (dht_await_pin_state(pin, 400, 0) != ESP_OK) {
        return 201; // Custom error: Phase B failed
    }
    
    // Phase C: Wait for sensor to release high
    if (dht_await_pin_state(pin, 400, 1) != ESP_OK) {
        return 202; // Custom error: Phase C failed
    }

    // Phase D: Wait for sensor to pull low to start data
    if (dht_await_pin_state(pin, 400, 0) != ESP_OK) {
        return 203; // Custom error: Phase D failed
    }

    // Read 40 bits
    for (int i = 0; i < 40; i++) {
        if (dht_await_pin_state(pin, 200, 1) != ESP_OK) return 204;
        
        int64_t start = esp_timer_get_time();
        if (dht_await_pin_state(pin, 200, 0) != ESP_OK) return 205;
        uint32_t high_duration = (uint32_t)(esp_timer_get_time() - start);

        uint8_t byte_idx = i / 8;
        data[byte_idx] <<= 1;
        if (high_duration > 40) data[byte_idx] |= 1;
    }

    return ESP_OK;
}

esp_err_t dht_read_data(dht_sensor_type_t sensor_type, gpio_num_t pin, int16_t *humidity, int16_t *temperature)
{
    uint8_t data[5] = { 0 };

    // 1. Reset & Start Signal
    gpio_config_t io_conf = {
        .intr_type = GPIO_INTR_DISABLE,
        .mode = GPIO_MODE_OUTPUT,
        .pin_bit_mask = (1ULL << pin),
        .pull_down_en = 0,
        .pull_up_en = 1,
    };
    gpio_config(&io_conf);
    
    gpio_set_level(pin, 1);
    vTaskDelay(pdMS_TO_TICKS(100));

    gpio_set_level(pin, 0);
    ets_delay_us(20000); // 20ms
    
    gpio_set_level(pin, 1);
    ets_delay_us(30);

    // 2. Switch to Input
    gpio_set_direction(pin, GPIO_MODE_INPUT);

    // 3. Critical Read
    PORT_ENTER_CRITICAL();
    esp_err_t result = dht_fetch_data_direct(pin, data);
    PORT_EXIT_CRITICAL();

    if (result != ESP_OK) return result;

    if (data[4] != ((data[0] + data[1] + data[2] + data[3]) & 0xFF)) return ESP_ERR_INVALID_CRC;

    if (humidity) {
        if (sensor_type == DHT_TYPE_DHT11) *humidity = data[0] * 10;
        else *humidity = ((data[0] << 8) + data[1]);
    }
    if (temperature) {
        if (sensor_type == DHT_TYPE_DHT11) *temperature = data[2] * 10;
        else {
            *temperature = (((data[2] & 0x7F) << 8) + data[3]);
            if (data[2] & 0x80) *temperature = -(*temperature);
        }
    }

    return ESP_OK;
}

esp_err_t dht_read_float_data(dht_sensor_type_t sensor_type, gpio_num_t pin, float *humidity, float *temperature)
{
    int16_t i_hum, i_temp;
    esp_err_t res = dht_read_data(sensor_type, pin, humidity ? &i_hum : NULL, temperature ? &i_temp : NULL);
    if (res != ESP_OK) return res;
    if (humidity) *humidity = i_hum / 10.0f;
    if (temperature) *temperature = i_temp / 10.0f;
    return ESP_OK;
}
