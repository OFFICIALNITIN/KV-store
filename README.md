# 🗄️ KV-Store

A lightweight, Redis-like **in-memory Key-Value store** written in Go. It communicates over raw TCP, supports TTL-based key expiration, and is fully configurable via a YAML file.

---

## ✨ Features

- ⚡ **Fast in-memory storage** using Go maps
- ⏱️ **TTL (Time-To-Live)** support — keys auto-expire
- 🔄 **Automatic janitor** that cleans up expired keys in the background
- 🌐 **TCP server** — connect with any TCP client (e.g. `netcat`, `telnet`)
- 🐳 **Docker-ready** — run anywhere with a single command
- ⚙️ **YAML config** — no recompile needed to change host/port/cleanup settings
- 🔒 **Concurrency-safe** using `sync.RWMutex`

---

## 📁 Project Structure

```
KV-store/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── server/
│   │   └── tcp.go           # TCP connection handler
│   └── store/
│       └── kvstore.go       # Core KV store logic
├── config.go                # Config loader (YAML → struct)
├── config.yaml              # Server configuration
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # One-command Docker setup
├── .gitignore
├── go.mod
└── go.sum
```

---

## ⚙️ Configuration

Edit `config.yaml` to configure the server — **no recompile needed**:

```yaml
server:
  host: 0.0.0.0   # Listen on all interfaces
  port: 6379       # TCP port

storage:
  cleanup_interval_seconds: 60   # How often to delete expired keys
  max_memory_mb: 512             # (informational, not enforced yet)
```

---

## 🚀 Running Locally (without Docker)

### Prerequisites
- [Go 1.24+](https://golang.org/dl/)

### Steps

```bash
# Clone the repository
git clone https://github.com/OFFICIALNITIN/KV-store.git
cd KV-store

# Download dependencies
go mod download

# Run the server
go run ./cmd/server/main.go
```

The server will start on `0.0.0.0:6379` by default.

---

## 🐳 Running with Docker

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### One-command start

```bash
docker compose up --build
```

> This builds the Go binary inside Docker and starts the server. No Go installation needed on your machine.

### Run in the background

```bash
docker compose up --build -d
```

### Stop the server

```bash
docker compose down
```

### Why it stays small 🐏
The `Dockerfile` uses a **multi-stage build**:
1. **Stage 1 (`builder`)** — Uses the full `golang:1.24-alpine` image to compile the binary
2. **Stage 2 (runtime)** — Copies only the compiled binary into a bare `alpine:latest` image (~5 MB)

The final Docker image contains **no Go toolchain** — just your server binary.

---

## 📡 TCP Protocol

Connect using `netcat`:

```bash
nc localhost 6379
```

### Supported Commands

| Command | Syntax | Description |
|---|---|---|
| `SET` | `SET <key> <value> [ttl_secs]` | Store a key-value pair, optionally with TTL |
| `GET` | `GET <key>` | Retrieve a value by key |
| `DELETE` | `DELETE <key>` | Remove a key |
| `EXIT` | `EXIT` | Close the connection |

### Examples

```
SET name nitin
OK

GET name
nitin

SET session abc123 30
OK

GET session
abc123

# (after 30 seconds)
GET session
(nil)

DELETE name
OK

EXIT
```

---

## 🏗️ Architecture

```
Client (netcat/telnet/your app)
        │  TCP
        ▼
  ┌─────────────┐
  │  tcp.go     │  ← One goroutine per connection
  │  (Handler)  │  ← Parses commands: SET / GET / DELETE
  └──────┬──────┘
         │
         ▼
  ┌─────────────┐
  │  kvstore.go │  ← Thread-safe map (sync.RWMutex)
  │  (Store)    │  ← TTL expiration logic
  └─────────────┘
         │
         ▼
  ┌─────────────┐
  │  Janitor    │  ← Background goroutine
  │  goroutine  │  ← Deletes expired keys every N seconds
  └─────────────┘
```

**Concurrency model:** Each TCP connection gets its own goroutine. The shared map is protected by `sync.RWMutex` — multiple clients can read simultaneously, but writes are exclusive.

---

## 🛠️ Development

```bash
# Run with race detector (detects concurrency bugs)
go run -race ./cmd/server/main.go

# Build binary manually
go build -o kv-server ./cmd/server/main.go

# Tidy dependencies
go mod tidy
```

---

## 📄 License

MIT — feel free to use, modify, and distribute.
