# Docker Deployment Guide

This guide covers how to deploy the in-memory storage service using Docker and Docker Compose.

## Prerequisites

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later

## Quick Start


1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd in-memory-storage
   ```

2. **Start the services:**
   ```bash
   make run-docker
   ```

3. **Verify the services are running:**
   ```bash
   docker-compose ps
   ```

4. **Access the services:**
   - **API Server**: http://localhost:8080
   - **Swagger UI**: http://localhost:8081


## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | Port for the HTTP server |
| `API_KEY` | `awesome-api-key` | API key for authentication |
