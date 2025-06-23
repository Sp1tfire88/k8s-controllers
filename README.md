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