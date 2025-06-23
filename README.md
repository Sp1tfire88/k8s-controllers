# k8s-controllers
🔐 GitHub Actions CI

| Job        | Purpose                                    |
| ---------- | ------------------------------------------ |
| `lint`     | Code formatting, `go vet`, `golangci-lint` |
| `test`     | Unit testing (`go test`)                   |
| `build`    | Compile project and validate binary        |
| `docker`   | Build Docker image                         |
| `security` | Trivy security scan of image               |

📄 Makefile Commands
| Command       | Description           |
| ------------- | --------------------- |
| `make build`  | Build Go binary       |
| `make run`    | Run the binary        |
| `make test`   | Run unit tests        |
| `make docker` | Build Docker image    |
| `make lint`   | Run golangci-lint     |
| `make tidy`   | Run `go mod tidy`     |
| `make clean`  | Clean build artifacts |

```
go tool cover -func=coverage/coverage.out
```
## ✅ Test Coverage Summary

| Function                | Coverage | Description                         |
|------------------------|----------|-------------------------------------|
| `initLogger`           | 88.9%    | Logging setup with Zerolog          |
| `init` (root.go)       | 81.8%    | CLI init: flags, Viper, bindings    |
| `Execute`              | 0.0%     | Cobra root command entrypoint       |
| `init` (server.go)     | 75.0%    | Server command init & flag binding  |
| `startFastHTTPServer`  | 0.0%     | Server start function (not tested)  |
| `logMiddleware`        | 0.0%     | Request logging middleware          |
| `homeHandler`          | 100.0%   | GET `/` handler                     |
| `postHandler`          | 100.0%   | POST `/post` handler                |
| `healthHandler`        | 100.0%   | GET `/health` handler               |
> 💡 See full [HTML coverage report](./coverage/coverage.html)
