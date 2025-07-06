# 🚀 API Gateway (Go)

A high-performance, plugin-based API Gateway built in Go — inspired by Kong and Apigee, but lightweight and open source.

---

## ✨ Features

- ⚡ Fast HTTP reverse proxy using Go’s standard library
- 🔌 Plugin system for dynamic middleware (auth, logging, rate limiting, etc.)
- 🗺️ Configurable routes via YAML
- 🔁 Route to multiple upstream services
- 📦 Designed for extensibility: add your own plugins easily

---

## 📦 Project Structure

api-gateway/
├── cmd/          # Entry point
├── config/       # YAML config loader
├── core/         # Plugin interface and registry
├── plugins/      # Built-in plugins (auth, logging, etc.)
├── proxy/        # Reverse proxy implementation
├── router/       # Route matching logic
├── test/         # Test helpers and integration tests
├── config.yaml   # Sample route configuration
├── Dockerfile
├── Makefile
└── go.mod

---

## ⚙️ Getting Started

### 🧱 Prerequisites

- Go >= 1.20
- Colima or Docker
- `make` (optional but recommended)

---

### 🏃 Run the Gateway

```bash
make run
```
Or manually:
```bash
go run ./cmd/main.go
```


⸻

🧪 Test It

With this example config:
```yaml
routes:
  - method: GET
    path: /hello
    upstream: https://httpbin.org
    plugins:
      - logging
```
Call it with:
```bash
curl http://localhost:8080/hello/get
```

⸻

🔌 Plugins

Plugins are registered in main.go and applied per route in config.

Example Plugins:
	•	logging: logs each request
	•	jwt-auth: blocks requests without Authorization header

You can add your own by implementing the Plugin interface.

⸻

🛠️ Development

Build
```bash
make build
```
Run Tests
```bash
make test
```
Build Docker Image
```bash
make docker-build
```

⸻

📚 Future Plans
	•	🔁 Retry/circuit breaker support
	•	🔐 mTLS and RBAC
	•	📈 Prometheus metrics & tracing
	•	🧬 gRPC support
	•	🌐 Admin API for live route changes
	•	🧩 Community plugin registry

⸻

🤝 Contributing

Contributions welcome! Submit issues, PRs, or plugin ideas. Let’s build an open-source gateway devs love.

⸻

📄 License

MIT — free to use, modify, and distribute.

---
