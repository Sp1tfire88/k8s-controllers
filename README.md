# k8s-controller
## Step 10 — Leader Election and Metrics for Controller Manager



$ go run main.go server
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-27T10:24:45Z","message":"🚀 Starting FastHTTP server on :8080"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-27T10:24:45Z","message":"🔧 Starting controller-runtime manager"}
2025-06-27T10:24:45Z    INFO    controller-runtime.metrics      Starting metrics server
2025-06-27T10:24:45Z    INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": ":9091", "secure": false}
2025-06-27T10:24:45Z    INFO    Starting EventSource    {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment", "source": "kind source: *v1.Deployment"}
2025-06-27T10:24:45Z    INFO    Starting Controller     {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment"}
2025-06-27T10:24:45Z    INFO    Starting workers        {"controller": "deployment", "controllerGroup": "apps", "controllerKind": "Deployment", "worker count": 1}
```


$ curl http://localhost:9091/metrics | head -n 10
```
# HELP certwatcher_read_certificate_errors_total Total number of certificate read errors
# TYPE certwatcher_read_certificate_errors_total counter
certwatcher_read_certificate_errors_total 0
# HELP certwatcher_read_certificate_total Total number of certificate reads
# TYPE certwatcher_read_certificate_total counter
certwatcher_read_certificate_total 0
# HELP controller_runtime_active_workers Number of currently used workers per controller
# TYPE controller_runtime_active_workers gauge
controller_runtime_active_workers{controller="deployment"} 0
# HELP controller_runtime_max_concurrent_reconciles Maximum number of concurrent reconciles per controller
# TYPE controller_runtime_max_concurrent_reconciles gauge
controller_runtime_max_concurrent_reconciles{controller="deployment"} 1
# HELP controller_runtime_reconcile_errors_total Total number of reconciliation errors per controller
# TYPE controller_runtime_reconcile_errors_total counter
controller_runtime_reconcile_errors_total{controller="deployment"} 0
```