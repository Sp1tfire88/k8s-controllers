# k8s-controller
## Step 9 — Integrating controller-runtime Deployment Controller

This step demonstrates how to integrate a Kubernetes controller using [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime), which watches for changes to `Deployment` resources and logs each reconcile event.

---

### ✅ What This Adds
- A `DeploymentReconciler` that reacts to `Deployment` create/update/delete events.
- A `controller-runtime` manager that runs alongside the FastHTTP server.
- Log messages for each reconciliation (Reconcile loop).

---

### 📂 Project Structure
```
├── cmd/
│   └── server.go       # Launches FastHTTP and controller-runtime manager
├── pkg/
│   └── controller/
│       └── controller.go  # Contains DeploymentReconciler logic
```

---

### 🔧 controller.go (Reconciler)
```go
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)
    logger.Info("🔁 Reconcile triggered", "name", req.Name, "namespace", req.Namespace)
    return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr manager.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&appsv1.Deployment{}).
        Complete(r)
}
```

---

### 🚀 server.go (Launch Manager)
```go
cfg := ctrl.GetConfigOrDie()
mgr, err := ctrl.NewManager(cfg, ctrl.Options{
    Scheme: scheme,
    Metrics: server.Options{
        BindAddress: ":8081",
    },
})

// Register controller
err = (&controller.DeploymentReconciler{
    Client: mgr.GetClient(),
    Scheme: mgr.GetScheme(),
}).SetupWithManager(mgr)

// Start controller-runtime
go func() {
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        log.Fatal().Err(err).Msg("controller-runtime manager exited")
    }
}()
```

---

### 🪵 Sample Log Output
```
2025-06-26T09:10:33Z    INFO    🔁 Reconcile triggered  {
  "controller": "deployment",
  "controllerGroup": "apps",
  "controllerKind": "Deployment",
  "Deployment": {"name":"test-nginx","namespace":"default"},
  "namespace": "default",
  "name": "test-nginx"
}
```

---

### 📦 How It Works
- The controller watches for events on `Deployment` resources.
- For each event, `Reconcile` is invoked with the object key (namespace + name).
- The controller fetches the latest object from the cache and logs the event.
- It is run via `controller-runtime` Manager with its own lifecycle and metrics server (port `:8081`).

---

### 🔁 Useful Commands for Testing
```bash
kubectl create deployment test-nginx --image=nginx
kubectl scale deployment test-nginx --replicas=2
kubectl delete deployment test-nginx
```
Each command triggers a reconcile event that is logged to the console.

---

