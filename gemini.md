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

The hardware component is based on the ESP32-C6 and a DHT22 sensor.

### Key Principles:
- **HTTP/JSON**: Data is transmitted to the server via standard POST requests in JSON format.
- **Secrets Management**: Sensitive data (WiFi SSID, password) is moved to a separate `secrets.h` file.
- **Modular Sensors**: The code is easily extensible to support multiple sensors.

## 4. Development Principles
- **Clean Code**: Clear folder structure and component naming.
- **Separation of Concerns**: Each layer has a strictly defined role.
- **Infrastructure as Code**: GitHub Actions (`.github/workflows/`) are present for build automation, linting, and deployment.

## 5. API Endpoints

### 5.1. Sensor Node API (`/api/v1`)
Endpoints used by ESP32 hardware nodes to transmit data.

- `POST /api/v1/measurements`: Submit new sensor readings.
    - **Payload**: `{ "device_id": "string", "temperature": float, "humidity": float, "timestamp": "RFC3339 string" }`
- `POST /api/v1/devices`: Register or update hardware device information.
    - **Payload**: `{ "id": "string", "mac": "string" }`

### 5.2. Admin Interface (`/admin`)
Server-side rendered (SSR) pages for data management and visualization.

- `GET /admin/`: Dashboard overview.
- `GET /admin/measurements`: Tabular view of historical data.
- `GET /admin/devices`: List of all registered sensor nodes.
- `GET /admin/devices/:id`: Detailed statistics and info for a specific node.
