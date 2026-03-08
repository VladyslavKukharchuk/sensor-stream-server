# Server-side Engineering Standards

This document defines the coding standards and architectural rules for the Go-based backend.

## 1. Environment Configuration
- **Rule**: `os.Getenv` and `godotenv` must ONLY be used in `main.go`.
- **Validation**: All required environment variables must be validated at startup. If a variable is missing, the application must terminate immediately using `log.Fatal()`.
- **Dependency Inversion**: Pass configuration values to controllers, services, and repositories via constructors (structs or arguments).

## 2. Layered Architecture
- **Controller**: Only handles HTTP-related logic (parsing body, status codes, rendering).
- **Service**: Contains business logic. Must be independent of Fiber/HTTP.
- **Repository**: Only handles data persistence (Firestore).

## 3. Formatting and Linting
- Follow standard Go conventions (`go fmt`).
- Avoid "magic numbers"; use constants for time durations, limits, etc.
- Pre-allocate slices with `make([]T, 0, length)` if the size is known to optimize memory usage.
- **Redundant Comments**: DO NOT add comments in obvious places (e.g., "// Controllers", "// Initialize", "// Routes"). Code should be self-documenting.
- **Code Style**: Prefer concise structures like `switch` statements over long `if-else` chains where it improves readability.
