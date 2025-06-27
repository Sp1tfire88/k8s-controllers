# k8s-controllers: Production-Ready Kubernetes Controller Example

**k8s-controllers** is a demonstration project of an advanced custom Kubernetes controller written in Go, using controller-runtime, FastHTTP REST API, cluster informer, full-featured CI, Helm/Kustomize manifests, leader election, metrics, and test coverage.

---

## 📦 Features

- **Go controller-runtime controller**
    - Reconciliation logic for Deployments (logs CRUD events).
    - Integrated FastHTTP REST API (`/deployments` endpoint returns Deployments from the informer's cache, not the live API).
- **Leader Election & Metrics**
    - CLI/config flags for leader election and metrics (`--enable-leader-election`, `--metrics-port`).
    - Prometheus endpoint exposed for metrics.
- **Helm & Kustomize** charts and CI validation.
- **Complete CI pipeline**: lint, test, build, docker, security scan, kustomize/helm validation, integration tests.
- **Configurable via CLI flags & `config.yaml`**.
- **Full test coverage**: unit, integration, and end-to-end tests (including leader election and API endpoints).
---

## 🚀 Quick Start

### 1. **Build & Run**
```bash
make build
./build/controller server --log-level trace --kubeconfig ~/.kube/config
```
Or using config file:
```
curl http://localhost:8080/deployments
# Example output: ["nginx-deployment","test-app", ...] (deployments from the informer's cache!)
```
| Parameter        | CLI Flag                 | config.yaml          | Default Value   |
| ---------------- | ------------------------ | -------------------- | --------------- |
| kubeconfig       | --kubeconfig             | kubeconfig           | \~/.kube/config |
| Logging          | --log-level              | log-level            | info            |
| REST port        | --port                   | port                 | 8080            |
| Metrics port     | --metrics-port           | metricsPort          | 8081            |
| Leader election  | --enable-leader-election | enableLeaderElection | false            |

🏗️ CI: Lint, Test, Security
The project uses GitHub Actions and provides a full pipeline:
* golangci-lint
* go test + coverage
* docker build
* Trivy security scan
* helm lint
* kustomize generate
* integration tests (leader election, metrics, API)