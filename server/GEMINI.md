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

## 4. Frontend & UI Standards (Server-Side Rendering)
The project uses Fiber's HTML template engine and follows a modular, component-based architecture for server-side rendering (SSR).

### 4.1. Directory Structure
- **`internal/views/layouts/`**: Base page structures (meta tags, global styles).
- **`internal/views/pages/`**: Route-specific content (e.g., dashboard, login).
- **`internal/views/components/`**: Atomic, reusable UI elements.

### 4.2. Component Architecture & Encapsulation
- **Principle**: Every reusable UI element (e.g., `device_card`, `header`) must be a standalone file in the `components/` directory.
- **Rule**: Components must be **self-contained**. HTML, CSS (via `<style>`), and JavaScript (via `<script>`) specific to a component should live within its file.
- **Rule**: Avoid adding component-specific styles to global CSS files like `styles.css`.

### 4.3. Layouts and Rendering
- **Principle**: Use Fiber's `{{embed}}` mechanism for layouts instead of complex block inheritance to prevent naming collisions.
- **Rule**: Controllers must explicitly specify the layout when calling `f.Render` (e.g., `f.Render("pages/device", data, "layouts/base")`).

### 4.4. Design Tokens & Styling
- **Principle**: Use CSS Variables for all visual styling (colors, spacing, shadows).
- **Rule**: All tokens must be defined in `public/css/variables.css`.
- **Rule**: Components should never use "magic" hex codes or hardcoded pixel values; always reference CSS variables (e.g., `var(--color-primary)`).

### 4.5. Clean JavaScript
- **Principle**: Prefer Vanilla JS for simple interactive elements (e.g., year updates, logout handling).
- **Rule**: Component-specific JS should be placed directly in the component file to maintain isolation.

## 5. Quality Assurance (Mandatory)
- **Rule**: Every code change in the `server/` directory MUST be followed by a successful `make server-lint` run.
- **Rule**: If `make server-lint` fails, the issue must be resolved before the task is considered complete.
- **Rule**: Use the `server-linter` skill to maintain quality standards.
