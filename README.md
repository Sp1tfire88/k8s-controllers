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
* cmd/namespaces.go

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

# added function for working with namespaces
$ go run main.go namespaces
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"warn","env":"dev","version":"v0.1.0","error":"no such flag -logtostderr","time":"2025-06-25T10:03:46Z","message":"Failed to set flag 'logtostderr'"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T10:03:46Z","message":"Found 5 namespace(s):"}
NAME             STATUS  AGE
default          Active  44m
kube-node-lease  Active  44m
kube-public      Active  44m
kube-system      Active  44m
test             Active  3m
```