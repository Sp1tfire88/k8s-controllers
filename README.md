# k8s-controller
## Deployment Informer with client-go

create go functions:
* informer.go

✅ Features:
* Watches 
* Deployment add/update/delete events.
* Supports both kubeconfig and in-cluster configuration.
* Reads configuration from config.yaml:
```
log-level: trace
port: 8080

kubeconfig: "/home/codespace/.kube/config"
informer:
  enabled: true

```
Start the informer

go run main.go server --log-level trace
```
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T13:29:16Z","message":"📦 Deployment ADDED"}
```
And when scaling:
```
{"level":"info","env":"dev","version":"v0.1.0","deployment":"test-nginx","from":1,"to":2,"time":"2025-06-25T13:29:46Z","message":"🔁 Deployment scaled"}
```
kubectl delete deployment test-nginx
```
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T13:31:21Z","message":"🗑️ Deployment DELETED"}
```