# Sensor Stream Project: System Manifest

This project is a complete IoT solution for collecting, storing, and visualizing sensor data (temperature and humidity) using ESP32-C6 hardware and a Go-based backend.

## 1. System Overview
The system consists of two main components:
- **`node/`**: ESP32 firmware for data collection and submission via HTTPS.
- **`server/`**: Go web server handling data persistence, authentication, and an administrative dashboard.

## 2. Global Architecture
The system follows a distributed architecture where hardware nodes act as data producers and the server acts as an orchestrator and data consumer.

### Component Interaction:
1. **Registration**: Hardware nodes register themselves using their MAC addresses.
2. **Collection**: Nodes read data from a DHT22 sensor every 5 minutes.
3. **Transmission**: Data is submitted to the server via secure REST API calls.
4. **Visualization**: Users monitor devices and historical data through a responsive web dashboard.

## 3. Communication API

### 3.1. Sensor Node API (`/api/v1`)
- `POST /api/v1/devices`: Register or update hardware device information.
    - **Payload**: `{ "id": "string", "mac": "string" }`
- `POST /api/v1/measurements`: Submit new sensor readings.
    - **Payload**: `{ "device_id": "string", "temperature": float, "humidity": float, "timestamp": "RFC3339 string" }`

### 3.2. Administrative Interface (`/admin`)
- `GET /admin/`: Dashboard showing the latest status of all registered devices.
- `GET /admin/devices/:id`: Detailed historical statistics with server-side aggregation:
    - `day`: 1-hour interval.
    - `week`: 6-hour interval.
    - `month`: 24-hour interval.
- `POST /admin/devices/:id`: Update device metadata (Name and Location).

### 3.3. Authentication API (`/auth`)
- `POST /auth/session`: Create a session cookie from a Firebase ID token.
- `DELETE /auth/session`: Destroy the current session.

## 4. Documentation Hierarchy
For detailed engineering standards and implementation rules, refer to the component-specific documentation:
- **`node/GEMINI.md`**: Firmware architecture and hardware development rules.
- **`server/GEMINI.md`**: Backend standards, UI components, and data processing logic.
