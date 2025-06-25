# k8s-controllers
List Kubernetes Deployments with client-go

intall dependencies
``
go get k8s.io/client-go@v0.30.0
go get k8s.io/apimachinery@v0.30.0
``
create cmd/list.go

create simple-deployment
```
sudo PATH=$PATH:/usr/sbin kubebuilder/bin/kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:latest
          securityContext:
            privileged: true
EOF
```
check go run main.go list --kubeconfig ./kubeconfig
```
 go run main.go list --kubeconfig ./kubeconfig
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"debug","env":"dev","version":"v0.1.0","time":"2025-06-25T08:52:06Z","message":"Using kubeconfig: ./kubeconfig"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T08:52:06Z","message":"Connected to cluster. Listing deployments..."}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-25T08:52:06Z","message":"Found 1 deployment(s):"}
📦 nginx-deployment
```