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
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
```
| Function              | Cover    | Comment                                          |
| --------------------- | -------- | ---------------------------------------------------- |
| `AddNewUser`          | ✅ 100%   | отлично                                              |
| `GetUsers`            | ✅ 100%   | отлично                                              |
| `initLogger`          | 🟡 89%   | почти, но не хватает ветвлений                       |
| `Execute`             | ❌ 0%     | **не вызывается напрямую** в тестах                  |
| `startFastHTTPServer` | ❌ 0%     | не тестируется (вызывает fasthttp.ListenAndServe)    |
| `logMiddleware`       | ❌ 0%     | не покрыт (нужно протестировать в контексте сервера) |
