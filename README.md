# Go Task Management RESTful API

This project is a lightweight, high-performance server application (RESTful API) designed to automate task tracking and management. Built with the Go programming language, the application focuses on idiomatic development patterns, minimal external dependencies, and predictable distributed system behavior.

## Architecture and Design Patterns

The core of the application relies on the foundational constraints of modern REST software architecture:

* **Separation of Concerns (Client-Server):** The backend is decoupled from any potential user interface, operates strictly over HTTP, and exchanges data using the unified JSON format.
* **Statelessness (Stateless):** Each incoming network request is processed independently. The server does not store client session contexts in RAM, ensuring straightforward horizontal scalability.
* **Layered System (Layered Architecture):** The codebase strictly separates responsibilities into network routing layers (Handlers), business logic elements, and low-level database operations (Repository / Store).
* **Safe Partial Updates (PATCH):** Data transfer models used for resource modification are built using Go pointers. This approach accurately distinguishes between implicit default values and the intentional absence of fields in a JSON payload, preventing accidental data overwrites in the database.

## Technology Stack

* **Runtime Environment:** Go 1.26
* **Database Management System:** PostgreSQL 15 (running inside an Alpine Linux container)
* **Database Driver and Tools:** `sqlx` (a general extension to standard `database/sql` for advanced struct mapping) and the official `pgx` / `lib/pq` drivers.
* **Containerization:** Docker, Docker Compose 3.8

## Getting Started

To spin up the project locally, you need the Go toolchain installed and the Docker Desktop daemon running on your host machine.

### 1. Database Deployment

Building the isolated environment and starting the PostgreSQL database is handled via a single command executed from the project root directory:

```bash
docker-compose up -d
```

During the initial container startup, the database automatically provisions a local Docker volume for persistent data storage and executes the initialization script located at `sql/init.sql`. This script sets up the necessary table structures, constraints, and inserts baseline test records.

### 2. Launching the Go Server

To compile and start the HTTP server, execute the following command:

```bash
go run main.go
```

Once the connection pool is initialized and verified via an internal system ping, the server will begin listening for incoming TCP traffic on the assigned port (port 9000 is configured by default to avoid common Windows system port constraints).

## API Specification and Examples

Below are standard examples for validating core endpoints using the cURL utility.

### Creating a New Task

* **Method:** POST
* **Path:** `http://localhost:9000/tasks/create`
* **Headers:** `Content-Type: application/json`

Execution via Windows PowerShell:
```powershell
curl.exe -X POST http://localhost:9000/tasks/create -H "Content-Type: application/json" -d '{\"title\": \"New task\", \"description\": \"Check the REST API\", \"completed\": true}'
```

Execution via Linux/macOS Bash:
```bash
curl -X POST http://localhost:9000/tasks/create \
  -H "Content-Type: application/json" \
  -d '{"title": "New task", "description": "Check the REST API", "completed": true}'
```

### Server Response (HTTP Status 201 Created)
```json
{
  "id": 6,
  "title": "New task",
  "description": "Check the REST API",
  "completed": true,
  "created_at": "2026-08-03T18:50:00Z",
  "updated_at": "2026-08-03T18:50:00Z"
}
```

## Connection Pool Resource Management

To prevent database performance degradation under heavy traffic, specific limits are applied to the connection pool lifecycle within the initialization code:
* Maximum active open connections: 20
* Maximum idle connections kept in reserve: 5
