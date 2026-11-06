# TaskForge

TaskForge is a distributed task-processing system built with Go and Redis. It allows clients to submit tasks through an HTTP API and have independent worker processes consume and execute them asynchronously.

---

### Tech Stack
Go – Backend application and worker processes.

Redis – Message broker and task queue.

Docker & Docker Compose – Running Redis locally.

net/http – HTTP API.

go-redis – Go client for Redis.

UUID – Unique task identifiers.

### Architecture
```text
       HTTP Client
            |
            v
       Go API Server
            |
            v
          Redis
       (Task Queue)
            |
            v
       Go Worker
            |
            v
      Task Processing
```
### Project Structure
```text
TaskForge/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── queue/
│   │   ├── queue.go
│   │   ├── queue_test.go
│   │   ├── redis.go
│   │   └── redis_test.go
│   └── task/
│       └── task.go
├── deployments/
│   └── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```
### Getting Started
Prerequisites:

Go installed

Docker Desktop installed and running

1. Start Redis
```bash
docker compose -f deployments/docker-compose.yml up -d
```

Verify Redis:

```bash
docker exec -it taskforge-redis redis-cli ping
```

Expected output:

```text
PONG
```

2. Start the API

In terminal 1:

```bash
go run ./cmd/api
```

The API runs on:

```text
http://localhost:8080
```
3. Start the Worker

In terminal 2:

```bash
go run ./cmd/worker
```

The worker waits for incoming tasks.

4. Submit a Task

In terminal 3:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name":"send-email"}'
```

The API returns the created task as JSON, and the worker consumes and processes it through Redis.

### Testing

Run all tests:

```bash
go test ./...
```

Run tests with the Go race detector:

```bash
go test -race ./...
```
### Current Features
HTTP endpoint for task creation (POST /tasks).

UUID-based task identification.

Redis-backed FIFO task queue.

Independent API and worker processes.

Blocking task consumption using Redis BRPOP.

Concurrent in-memory queue implementation with mutex protection.

Unit tests for queue operations.

### Future Improvements
Reliable task delivery and acknowledgements.

Automatic retries and failed-task handling.

PostgreSQL persistence.

Multiple worker instances.

Task status tracking.

Graceful shutdown and improved error handling.

---
TaskForge – Learning and building distributed systems, one task at a time.