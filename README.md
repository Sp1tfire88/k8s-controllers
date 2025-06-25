# k8s-controllers
## Step6+: Add create/delete command

intall dependencies:
```
go get k8s.io/client-go@v0.30.0
go get k8s.io/apimachinery@v0.30.0
```
create go functions:
* cmd/list.go
* cmd/create.go
* cmd/delete.go

$ go run main.go create --kubeconfig ~/.kube/config --name nginx-app --image nginx:latest --replicas 2
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T09:10:59Z","message":"✅ Deployment \"nginx-app\" created"}
```
$ go run main.go list --kubeconfig ~/.kube/config
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"debug","env":"dev","version":"v0.1.0","time":"2025-06-25T09:11:03Z","message":"Using kubeconfig: /home/codespace/.kube/config"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T09:11:03Z","message":"Connected to cluster. Listing deployments..."}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T09:11:03Z","message":"Found 1 deployment(s):"}
📦 nginx-app
```
$ go run main.go delete --kubeconfig ~/.kube/config --name nginx-app
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T09:11:07Z","message":"🗑️ Deployment \"nginx-app\" deleted"}
```