# Node Engineering Standards (ESP32-C6)

This document defines the architecture and development principles for the hardware sensor nodes.

## 1. Architecture Overview
The firmware is based on the ESP-IDF framework and follows a **Modular Component-Based Architecture**.

### 1.1. Core Components (`node/components/`)
Each major functionality is isolated into a standalone ESP-IDF component:
- **`dht/`**: Low-level driver for the DHT22 sensor (bit-banging protocol).
- **`app_wifi/`**: Connection management, event handling, and auto-reconnect.
- **`app_http/`**: Network communication, device registration, and data submission via HTTPS.
- **`app_storage/`**: NVS (Non-Volatile Storage) management for persisting `device_id`.
- **`app_time/`**: SNTP synchronization for accurate data timestamps.

### 1.2. Orchestration (`main/main.c`)
The main orchestrator handles:
1. System-wide initialization (NVS, WiFi).
2. The startup sequence: WiFi -> Time Sync -> Storage Check -> Device Registration.
3. Spawning the **`sensor_task`** (FreeRTOS task) for the main data collection loop.

## 2. Development Principles
- **Encapsulation**: Components must communicate via header interfaces (`.h`). Direct access to low-level structures from `main.c` is discouraged.
- **Naming Convention**: Application-specific components must be prefixed with `app_` to distinguish them from standard ESP-IDF libraries.
- **Robustness**: The registration logic must be idempotent. If a `device_id` is missing, the node automatically registers using its MAC address.
- **Security**: All API communication must use HTTPS with the ESP-IDF CRT Bundle for certificate validation.

## 3. Configuration
- **Secrets**: Environment-specific variables (SSID, URLs) are managed via `secrets.h` (based on `secrets.template.h`).
- **NVS**: Used to store and persist the unique `device_id` provided by the server during registration.
