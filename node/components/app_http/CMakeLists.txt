idf_component_register(SRCS "app_http.c"
                    INCLUDE_DIRS "include"
                    PRIV_REQUIRES esp_http_client json mbedtls esp-tls)
