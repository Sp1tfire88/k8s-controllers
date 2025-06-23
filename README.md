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
| `AddNewUser`          | ✅ 100%   | great                                              |
| `GetUsers`            | ✅ 100%   | great                                              |
| `initLogger`          | 🟡 89%    | almost, but not enough branches                      |
| `Execute`             | ❌ 0%     | not called directly in tests                         |
| `startFastHTTPServer` | ❌ 0%     | not tested (calls fasthttp.ListenAndServe)    |
| `logMiddleware`       | ❌ 0%     | not covered (needs to be tested in server context) |
