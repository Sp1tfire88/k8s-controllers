# k8s-controller
## /deployments JSON API Endpoint

$ go run main.go server --log-level trace

$ go run main.go create --name test-nginx --image nginx:1.25.2 --replicas 1
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"warn","env":"dev","version":"v0.1.0","error":"no such flag -logtostderr","time":"2025-06-25T13:48:46Z","message":"Failed to set flag 'logtostderr'"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T13:48:46Z","message":"✅ Deployment \"test-nginx\" created"}
```
$ go run main.go create --name redis --image redis --replicas 1
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"warn","env":"dev","version":"v0.1.0","error":"no such flag -logtostderr","time":"2025-06-25T13:49:06Z","message":"Failed to set flag 'logtostderr'"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T13:49:06Z","message":"✅ Deployment \"redis\" created"}
```
$ curl http://localhost:8080/deployments
```
["test-nginx","redis"]
```