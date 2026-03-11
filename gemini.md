# Architecture and Principles of the Sensor Stream Project

This project is a system for collecting, storing, and displaying sensor data (temperature and humidity). It consists of two main parts: a hardware node and a server-side component.

## 1. General Structure

The project is divided into two main directories:
- `node/`: Firmware for the ESP32 microcontroller (Arduino/C++).
- `server/`: Go-based web server for data processing and visualization.

## 2. Server Architecture (Go)

The server is built on **Layered Architecture** principles, ensuring a clear **Separation of Concerns**:

### Core Layers:
- **Controller (`internal/controller/`)**: Handles incoming HTTP requests, validates input data, and returns responses (JSON or HTML via templates).
- **Service (`internal/service/`)**: Contains the application's business logic. Uses interfaces to interact with repositories, facilitating easier testing.
- **Repository (`internal/repository/`)**: Responsible for data persistence and retrieval. Abstracting database operations from the rest of the code.
- **Model (`internal/model/`)**: Defines data structures used across all application layers.
- **Routes (`internal/routes/`)**: Defines endpoints and binds them to their respective controllers.
- **Views (`internal/views/`)**: HTML templates for server-side page rendering.

### Tech Stack and Principles:
- **Framework**: [Fiber](https://gofiber.io/) — a high-performance web framework inspired by Express.js.
- **Database**: [Google Cloud Firestore](https://cloud.google.com/firestore) — a NoSQL database for storing measurements and device information.
- **Logging**: [zerolog](https://github.com/rs/zerolog) for structured logging.
- **Dependency Injection**: Manual dependency injection is applied in the `main.go` file.
- **Environment Config**: Uses `.env` files for configuration (via `godotenv`).
- **Containerization**: Includes a `Dockerfile` for deployment in Docker and Google Cloud Run.

## 3. Node Architecture (ESP32)

The hardware component is based on the ESP32-C6 and a DHT22 sensor. It follows a **Modular Component-Based Architecture** to ensure separation of concerns and maintainability.

### 3.1. Components Structure
Each major functionality is isolated into a standalone ESP-IDF component within the `node/components/` directory:

- **`dht/`**: Low-level driver for the DHT22 sensor. Handles microsecond-accurate timing and bit-banging protocol.
- **`app_wifi/`**: Manages WiFi connectivity, event handling (auto-reconnect), and IP acquisition.
- **`app_http/`**: Handles all network communication with the server, including device registration and measurement submission via HTTPS.
- **`app_storage/`**: Manages non-volatile storage (SPIFFS) for persisting the `device_id`.
- **`app_time/`**: Handles SNTP synchronization to ensure accurate timestamps for sensor data.

### 3.2. Orchestration (`main.c`)
The `main/main.c` file acts as an **orchestrator**. Its responsibilities are:
1. Initializing system-wide resources (NVS).
2. Coordinating the startup sequence (WiFi -> Time -> Storage -> Registration).
3. Spawning the **`sensor_task`** (a FreeRTOS task) that runs the main infinite loop for data collection.

### 3.3. Key Principles:
- **Encapsulation**: Components communicate via clean header interfaces (`.h` files). `main.c` should not access low-level WiFi or HTTP structures directly.
- **Modular**: The code is easily extensible to support multiple sensors.
- **Robustness**: The registration logic is idempotent. If a `device_id` is missing, the node automatically registers itself using its MAC address.
- **Security**: HTTPS is used for all API calls, utilizing the ESP-IDF CRT Bundle for certificate validation. Sensitive data is stored in `secrets.h`.

### 3.4. Development Guidelines:
- **Adding Features**: Create a new component in `components/` if the feature is reusable or complex.
- **Configuration**: Always use `secrets.h` for environment-specific variables (SSID, URLs).
- **Naming**: Prefix application-specific components with `app_` to distinguish them from standard ESP-IDF libraries.

## 4. Development Principles
- **Clean Code**: Clear folder structure and component naming.
- **Separation of Concerns**: Each layer has a strictly defined role.
- **Infrastructure as Code**: GitHub Actions (`.github/workflows/`) are present for build automation, linting, and deployment.

## 5. API Endpoints

### 5.1. Authentication API (`/auth`)
Endpoints for user session management via Firebase ID tokens.

- `POST /auth/session`: Create a session cookie from a Firebase ID token.
    - **Payload**: `{ "idToken": "string" }`
- `DELETE /auth/session`: Destroy the current session cookie (Logout).

### 5.2. Sensor Node API (`/api/v1`)
Endpoints used by ESP32 hardware nodes to transmit data.

- `POST /api/v1/measurements`: Submit new sensor readings.
    - **Payload**: `{ "device_id": "string", "temperature": float, "humidity": float, "timestamp": "RFC3339 string" }`
- `POST /api/v1/devices`: Register or update hardware device information.
    - **Payload**: `{ "id": "string", "mac": "string" }`

### 5.2. Admin Interface (`/admin`)
Server-side rendered (SSR) pages for data management and visualization.

- `GET /admin/`: Main dashboard showing interactive cards for all registered devices with their latest sensor readings and online status.
- `GET /admin/devices/:id`: Detailed statistics and interactive charts for a specific node.
    - **Query Params**: `?period=day|week|month` (default: `day`).
    - **Data Processing**: Automatically applies server-side aggregation (averaging) based on the period:
        - `day`: 1-hour interval.
        - `week`: 6-hour interval.
        - `month`: 24-hour interval.
- `POST /admin/devices/:id`: Update device metadata (Friendly Name and Location).
    - **Form Data**: `name` (string), `location` (string).

### 5.3. Static Content (`/`)
- Serves frontend assets (CSS, JS, images) from the `./public` directory.
- ApexCharts is loaded via CDN for interactive data visualization.
