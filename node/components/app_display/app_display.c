#include "app_display.h"

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "driver/i2c_master.h"
#include "esp_lcd_io_i2c.h"
#include "esp_lcd_panel_io.h"
#include "esp_lcd_panel_ops.h"
#include "esp_lcd_panel_vendor.h"

#define DISPLAY_WIDTH 128
#define DISPLAY_HEIGHT 64
#define DISPLAY_I2C_PORT I2C_NUM_0
#define DISPLAY_I2C_ADDRESS 0x3C
#define DISPLAY_SDA_GPIO 2
#define DISPLAY_SCL_GPIO 1
#define DISPLAY_I2C_FREQUENCY_HZ 100000
#define DISPLAY_FRAMEBUFFER_SIZE (DISPLAY_WIDTH * DISPLAY_HEIGHT / 8)
#define FONT_WIDTH 5
#define FONT_HEIGHT 7
#define FONT_COLUMN_SPACING 1

static i2c_master_bus_handle_t s_i2c_bus;
static esp_lcd_panel_io_handle_t s_panel_io;
static esp_lcd_panel_handle_t s_panel;
static uint8_t s_framebuffer[DISPLAY_FRAMEBUFFER_SIZE];
static bool s_initialized;

static const uint8_t *glyph_rows(char character)
{
    static const uint8_t blank[FONT_HEIGHT] = {0};
    static const uint8_t glyph_a[FONT_HEIGHT] = {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11};
    static const uint8_t glyph_c[FONT_HEIGHT] = {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E};
    static const uint8_t glyph_d[FONT_HEIGHT] = {0x1C, 0x12, 0x11, 0x11, 0x11, 0x12, 0x1C};
    static const uint8_t glyph_e[FONT_HEIGHT] = {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F};
    static const uint8_t glyph_g[FONT_HEIGHT] = {0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0F};
    static const uint8_t glyph_h[FONT_HEIGHT] = {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11};
    static const uint8_t glyph_i[FONT_HEIGHT] = {0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E};
    static const uint8_t glyph_m[FONT_HEIGHT] = {0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11};
    static const uint8_t glyph_n[FONT_HEIGHT] = {0x11, 0x19, 0x19, 0x15, 0x13, 0x13, 0x11};
    static const uint8_t glyph_o[FONT_HEIGHT] = {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E};
    static const uint8_t glyph_p[FONT_HEIGHT] = {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10};
    static const uint8_t glyph_r[FONT_HEIGHT] = {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11};
    static const uint8_t glyph_s[FONT_HEIGHT] = {0x0F, 0x10, 0x10, 0x0E, 0x01, 0x01, 0x1E};
    static const uint8_t glyph_t[FONT_HEIGHT] = {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04};
    static const uint8_t glyph_u[FONT_HEIGHT] = {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E};
    static const uint8_t glyph_0[FONT_HEIGHT] = {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E};
    static const uint8_t glyph_1[FONT_HEIGHT] = {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E};
    static const uint8_t glyph_2[FONT_HEIGHT] = {0x0E, 0x11, 0x01, 0x02, 0x04, 0x08, 0x1F};
    static const uint8_t glyph_3[FONT_HEIGHT] = {0x1E, 0x01, 0x01, 0x0E, 0x01, 0x01, 0x1E};
    static const uint8_t glyph_4[FONT_HEIGHT] = {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02};
    static const uint8_t glyph_5[FONT_HEIGHT] = {0x1F, 0x10, 0x10, 0x1E, 0x01, 0x01, 0x1E};
    static const uint8_t glyph_6[FONT_HEIGHT] = {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E};
    static const uint8_t glyph_7[FONT_HEIGHT] = {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08};
    static const uint8_t glyph_8[FONT_HEIGHT] = {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E};
    static const uint8_t glyph_9[FONT_HEIGHT] = {0x0E, 0x11, 0x11, 0x0F, 0x01, 0x02, 0x0C};
    static const uint8_t glyph_colon[FONT_HEIGHT] = {0x00, 0x04, 0x04, 0x00, 0x04, 0x04, 0x00};
    static const uint8_t glyph_dot[FONT_HEIGHT] = {0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x0C};
    static const uint8_t glyph_minus[FONT_HEIGHT] = {0x00, 0x00, 0x00, 0x1F, 0x00, 0x00, 0x00};
    static const uint8_t glyph_percent[FONT_HEIGHT] = {0x18, 0x19, 0x02, 0x04, 0x08, 0x13, 0x03};

    switch (character) {
    case 'A': return glyph_a;
    case 'C': return glyph_c;
    case 'D': return glyph_d;
    case 'E': return glyph_e;
    case 'G': return glyph_g;
    case 'H': return glyph_h;
    case 'I': return glyph_i;
    case 'M': return glyph_m;
    case 'N': return glyph_n;
    case 'O': return glyph_o;
    case 'P': return glyph_p;
    case 'R': return glyph_r;
    case 'S': return glyph_s;
    case 'T': return glyph_t;
    case 'U': return glyph_u;
    case '0': return glyph_0;
    case '1': return glyph_1;
    case '2': return glyph_2;
    case '3': return glyph_3;
    case '4': return glyph_4;
    case '5': return glyph_5;
    case '6': return glyph_6;
    case '7': return glyph_7;
    case '8': return glyph_8;
    case '9': return glyph_9;
    case ':': return glyph_colon;
    case '.': return glyph_dot;
    case '-': return glyph_minus;
    case '%': return glyph_percent;
    case ' ': return blank;
    default: return blank;
    }
}

static void draw_text(const char *text, int x, int y, int scale)
{
    for (const char *character = text; *character != '\0'; ++character) {
        const uint8_t *rows = glyph_rows(*character);
        for (int row = 0; row < FONT_HEIGHT; ++row) {
            for (int column = 0; column < FONT_WIDTH; ++column) {
                if ((rows[row] & (1U << (FONT_WIDTH - column - 1))) == 0) {
                    continue;
                }

                for (int dy = 0; dy < scale; ++dy) {
                    for (int dx = 0; dx < scale; ++dx) {
                        int pixel_x = x + column * scale + dx;
                        int pixel_y = y + row * scale + dy;
                        size_t index = (size_t)(pixel_y / 8) * DISPLAY_WIDTH + pixel_x;
                        s_framebuffer[index] |= (uint8_t)(1U << (pixel_y % 8));
                    }
                }
            }
        }
        x += (FONT_WIDTH + FONT_COLUMN_SPACING) * scale;
    }
}

static esp_err_t render_frame(void)
{
    return esp_lcd_panel_draw_bitmap(s_panel, 0, 0, DISPLAY_WIDTH, DISPLAY_HEIGHT, s_framebuffer);
}

static void clear_frame(void)
{
    memset(s_framebuffer, 0, sizeof(s_framebuffer));
}

esp_err_t app_display_init(void)
{
    const i2c_master_bus_config_t bus_config = {
        .clk_source = I2C_CLK_SRC_DEFAULT,
        .glitch_ignore_cnt = 7,
        .i2c_port = DISPLAY_I2C_PORT,
        .sda_io_num = DISPLAY_SDA_GPIO,
        .scl_io_num = DISPLAY_SCL_GPIO,
        .flags.enable_internal_pullup = true,
    };
    esp_err_t err = i2c_new_master_bus(&bus_config, &s_i2c_bus);
    if (err != ESP_OK) {
        return err;
    }

    const esp_lcd_panel_io_i2c_config_t io_config = {
        .dev_addr = DISPLAY_I2C_ADDRESS,
        .scl_speed_hz = DISPLAY_I2C_FREQUENCY_HZ,
        .control_phase_bytes = 1,
        .dc_bit_offset = 6,
        .lcd_cmd_bits = 8,
        .lcd_param_bits = 8,
        .transaction_timeout_ms = 1000,
    };
    err = esp_lcd_new_panel_io_i2c(s_i2c_bus, &io_config, &s_panel_io);
    if (err != ESP_OK) {
        return err;
    }

    const esp_lcd_panel_dev_config_t panel_config = {
        .bits_per_pixel = 1,
        .reset_gpio_num = -1,
    };
    err = esp_lcd_new_panel_ssd1306(s_panel_io, &panel_config, &s_panel);
    if (err != ESP_OK) {
        return err;
    }

    err = esp_lcd_panel_reset(s_panel);
    if (err != ESP_OK) {
        return err;
    }
    err = esp_lcd_panel_init(s_panel);
    if (err != ESP_OK) {
        return err;
    }
    err = esp_lcd_panel_disp_on_off(s_panel, true);
    if (err != ESP_OK) {
        return err;
    }

    s_initialized = true;
    return app_display_show_startup();
}

esp_err_t app_display_show_startup(void)
{
    if (!s_initialized) {
        return ESP_ERR_INVALID_STATE;
    }

    clear_frame();
    draw_text("STARTING", 16, 26, 2);
    return render_frame();
}

esp_err_t app_display_show_reading(float temperature, float humidity)
{
    if (!s_initialized) {
        return ESP_ERR_INVALID_STATE;
    }

    char temperature_text[16];
    char humidity_text[16];
    snprintf(temperature_text, sizeof(temperature_text), "T: %.1fC", (double)temperature);
    snprintf(humidity_text, sizeof(humidity_text), "H: %.1f%%", (double)humidity);

    clear_frame();
    draw_text("SENSOR DATA", 4, 4, 1);
    draw_text(temperature_text, 4, 19, 2);
    draw_text(humidity_text, 4, 41, 2);
    return render_frame();
}

esp_err_t app_display_show_sensor_error(void)
{
    if (!s_initialized) {
        return ESP_ERR_INVALID_STATE;
    }

    clear_frame();
    draw_text("DHT ERROR", 10, 26, 2);
    return render_frame();
}
