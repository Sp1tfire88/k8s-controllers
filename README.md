# k8s-controller
## Deployment Informer with client-go

create go functions:
* informer.go

$ go run main.go server --kubeconfig ~/.kube/config --log-level trace
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T12:51:30Z","message":"Starting FastHTTP server on :8080"}
{"level":"warn","env":"dev","version":"v0.1.0","error":"no such flag -logtostderr","time":"2025-06-25T12:51:30Z","message":"Failed to set flag 'logtostderr'"}
{"level":"trace","env":"dev","version":"v0.1.0","kubeconfig":"/home/codespace/.kube/config","time":"2025-06-25T12:51:30Z","message":"Using external kubeconfig"}
{"level":"trace","env":"dev","version":"v0.1.0","namespace":"default","time":"2025-06-25T12:51:30Z","message":"Creating informer factory"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T12:51:30Z","message":"🚀 Starting deployment informer"}
```
$ go run main.go create --name test-nginx --image nginx:1.25.2 --replicas 1
```
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T12:51:41Z","message":"📦 Deployment ADDED"}
```
$ kubectl set image deployment/test-nginx test-nginx=nginx:1.25.3
```
{"level":"trace","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T12:52:39Z","message":"✏️ Deployment UPDATED (no replica change)"}
```
$ kubectl scale deployment test-nginx --replicas=4
```
{"level":"info","env":"dev","version":"v0.1.0","deployment":"test-nginx","old":1,"new":4,"time":"2025-06-25T12:53:11Z","message":"🔁 Replicas count changed"}
```
$ kubectl delete deployment test-nginx
```
{"level":"info","env":"dev","version":"v0.1.0","deployment":"test-nginx","time":"2025-06-25T12:54:23Z","message":"✅ Confirmed deletion from cache"}
```