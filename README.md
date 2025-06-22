# Coding test for Go developers: In-memory data structure store

Need to implement simple in-memory data structure store. See Redis for example.

Required data strutures:
- Strings
- Lists

Required operations:
- Get
- Set
- Update
- Remove
- Push for lists
- Pop for lists

Required features:
- Keys with a limited TTL
- Go client API library
- HTTP REST API

Add unit tests for Go API and integration tests for REST API (without full coverage, just for example).
Provide REST API specs with examples, client library API docs and deployment docs (for Docker).

Optional features:
- Data persistence
- Perfomance tests
- Authentication

## Quick Start

-  **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd in-memory-storage
   ```
- **Install dependencies**
   ```bash
   make mod
   ```

### Using Docker (Recommended)

1. **Start the services:**
   ```bash
   make run-docker
   ```

2. **Access the services:**
   - **API Server**: http://localhost:8080
   - **Swagger UI**: http://localhost:8081

### Using Go directly

1. **Build and run:**
   ```bash
   make run
   ```

## Documentation

- **[Storage Library API](docs/storage_api.md)** - Complete API documentation for the storage library
- **[OpenAPI Specification](docs/openapi.yaml)** - REST API specification
- **[Postman Collection](docs/postman_collection.json)** - Ready-to-use API testing collection
- **[Docker Deployment Guide](docs/docker_deployment.md)** - Comprehensive Docker deployment documentation

## Testing

You can find unit testing for most of the code and a couple of examples of E2E testing in the `http_test.go` file.

To run all the tests:
```bash
make test
```

## Linter
You can run the default golangci-lint by running:
```bash
make lint
```

## Features Implemented

✅ **Core Data Structures**
- String storage with TTL support
- List storage with TTL support
- Generic list implementation for any data type. This could be expanded to support different data types.

✅ **Required Operations**
- Get, Set, Update, Remove for strings and lists
- Push and Pop operations for lists (FIFO)
- Thread-safe operations with locking

✅ **HTTP REST API**
- Complete REST API with authentication
- JSON request/response format
- Comprehensive error handling and logging
- OpenAPI 3.0 specification

✅ **Testing**
- Unit tests for storage library
- Integration tests for HTTP endpoints
- End-to-end tests for critical paths

✅ **Optional Features**
- API key authentication


## API Authentication

All API endpoints require authentication using an API key. Include the key in the Authorization header:

```
Authorization: Bearer awesome-api-key
```


## Project Structure

```
├── cmd/server/           # Main application entry point
├── internal/             # Internal application code
│   ├── app/             # Application setup and configuration
│   ├── http/            # HTTP server and middleware
│   ├── strings/         # String controller and models
│   └── lists/           # List controller and models
├── storage/             # Core storage library
├── docs/                # Documentation
│   ├── storage_api.md   # Storage library documentation
│   ├── openapi.yaml     # API specification
│   ├── postman_collection.json # API testing collection
│   └── docker_deployment.md # Docker deployment guide
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Multi-service deployment
└── README.md           # This file
```
