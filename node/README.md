# sensor-stream-sensor-node

This project demonstrates reading temperature and humidity from a **DHT22 (AM2302)** sensor using the **ESP32-C6-Zero** board and send data to web server.

---

## 📦 Components

- ESP32-C6-Zero
- DHT22 (AM2302)
- 4.7-10 kΩ resistor (for pull-up on DATA)
- Jumper wires or breadboard

---

## 🔌 Wiring

![wiring.png](wiring.png)
- 3.3V --> VCC DHT22
- GND --> GND DHT22
- GPIO4 --> DATA DHT22 (through 4.7-10kΩ resistor to 3.3V)

## ⚡ Features

- Reads temperature and humidity from DHT22 sensor.
- Sends JSON payload to server **POST /api/v1/measurements**:
    ```json
    {
      "temperature": 23.5,
      "humidity": 60,
       "timestamp": "2025-11-11T15:45:00Z"
    }
    ```
- Easy to extend for multiple ESP32 sensors.

### 🔧 Setup

Install the PlatformIO IDE extension in VS Code and open the `node` directory.
Run **PlatformIO: Build** for the `esp32-c6` environment. PlatformIO installs the
pinned Espressif32 7.1.3 platform (ESP-IDF 6.1) and managed `cjson` component
automatically; no separate ESP-IDF installation or Arduino libraries are needed.

The project uses `sdkconfig.defaults` and `partitions.csv` for its 8 MB flash
configuration and a 2 MB application partition. Before flashing, configure
`main/secrets.h` from `main/secrets.template.h`.

For this personal project, TLS server-certificate verification is disabled so
the node can connect even if its CA bundle or clock is not set up. HTTPS traffic
remains encrypted, but the node cannot verify the server's identity and is
vulnerable to man-in-the-middle attacks. Do not use this setting for sensitive
data or production deployments.

### Build, flash, and monitor

- Build with **PlatformIO: Build**.
- Connect the board and flash with **PlatformIO: Upload**. Select the correct
  serial port if PlatformIO does not detect it automatically.
- Open **PlatformIO: Monitor** at 115200 baud. Press `Ctrl+]` to exit.

### ⚡ Notes

- Ensure the pull-up resistor is installed between DATA and VCC; otherwise, readings may fail.
- If nothing appears in the Serial Monitor:
  - Check the correct COM port
  - Press RST on the board
  - Ensure baud rate = 115200
