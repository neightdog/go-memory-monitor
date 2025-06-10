# Go Memory, Disk & CPU Monitor (Microservices + RabbitMQ)

This project demonstrates a microservices architecture in Go, using RabbitMQ for pub/sub communication and a web UI for real-time system monitoring.

## Architecture

- **memory-service:** Monitors system memory usage and publishes updates to RabbitMQ.
- **disk-service:** Monitors disk usage and publishes updates to RabbitMQ.
- **cpu-service:** Monitors CPU usage and publishes updates to RabbitMQ.
- **web-ui:** Subscribes to RabbitMQ, exposes REST endpoints, and serves a live dashboard.

## Features

- Real-time system monitoring (memory, disk, and CPU)
- Decoupled microservices using RabbitMQ (fanout exchange)
- REST API and web dashboard
- Cross-platform (macOS, Linux)
- Easily extensible (add more metrics or consumers)

## Running Locally

1. **Start RabbitMQ** (Docker recommended):
   ```
   docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
   ```
2. **Run each service in a separate terminal:**
   ```
   go run ./cmd/memory-service
   go run ./cmd/disk-service
   go run ./cmd/cpu-service
   go run ./cmd/web-ui
   ```
3. **Open [http://localhost:8080](http://localhost:8080) in your browser.**

## Running with Docker

1. Build and start all services:
   ```
   docker-compose up --build
   ```
2. Visit [http://localhost:8080](http://localhost:8080) for the dashboard.
3. Access RabbitMQ management at [http://localhost:15672](http://localhost:15672) (guest/guest).

## Configuration

- RabbitMQ URL, disk path, and other settings can be set via environment variables.

## Why RabbitMQ?

- Demonstrates pub/sub, decoupling, and scalable event-driven architecture.
- Easy to add new consumers (e.g., alerting, logging) without changing publishers.

---

**Now supports monitoring of memory, disk, and CPU usage!**