# k8s-controllers
## Deployment Informer with client-go

create go functions:
* informer.go

$ go run main.go server --kubeconfig ~/.kube/config --log-level trace
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T10:52:41Z","message":"Starting FastHTTP server on :8080"}
{"level":"warn","env":"dev","version":"v0.1.0","error":"no such flag -logtostderr","time":"2025-06-25T10:52:41Z","message":"Failed to set flag 'logtostderr'"}
{"level":"trace","env":"dev","version":"v0.1.0","kubeconfig":"/home/codespace/.kube/config","time":"2025-06-25T10:52:41Z","message":"Using external kubeconfig"}
{"level":"trace","env":"dev","version":"v0.1.0","namespace":"default","time":"2025-06-25T10:52:41Z","message":"Creating informer factory"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T10:52:41Z","message":"🚀 Starting deployment informer"}
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T10:52:41Z","message":"📦 Deployment ADDED"}
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"nginx-app","time":"2025-06-25T10:52:49Z","message":"📦 Deployment ADDED"}
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"nginx-app","time":"2025-06-25T10:53:02Z","message":"🗑️ Deployment DELETED"}
```