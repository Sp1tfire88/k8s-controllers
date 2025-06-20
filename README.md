# k8s-controllers
Zerolog

## build
```go build -o controller
```
## run
$ ./controller 
```
13:06:55 INF This is an info log env=dev version=v0.1.0
13:06:55 WRN This is a warn log env=dev version=v0.1.0
13:06:55 ERR This is an error log env=dev version=v0.1.0
Welcome to k8s-controller-tutorial CLI!
```

## run with debug output

./controller --log-level=trace
``` 
13:07:14 TRC This is a trace log env=dev version=v0.1.0
13:07:14 DBG This is a debug log env=dev version=v0.1.0
13:07:14 INF This is an info log env=dev version=v0.1.0
13:07:14 WRN This is a warn log env=dev version=v0.1.0
13:07:14 ERR This is an error log env=dev version=v0.1.0
Welcome to k8s-controller-tutorial CLI!
```
## log entry in the output file
logs\app.log
```
{"level":"trace","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:34Z","message":"This is a trace log"}
{"level":"debug","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:34Z","message":"This is a debug log"}
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:34Z","message":"This is an info log"}
{"level":"warn","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:34Z","message":"This is a warn log"}
{"level":"error","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:34Z","message":"This is an error log"}
{"level":"error","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:49Z","message":"This is an error log"}
{"level":"error","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:57Z","message":"This is an error log"}
{"level":"error","env":"dev","version":"v0.1.0","time":"2025-06-20T13:15:58Z","message":"This is an error log"}
```

# Use Viper to add env vars*
```
go get github.com/spf13/viper
```

```
$ LOG_LEVEL=info ./controller
13:29:19 INF This is an info log env=dev version=v0.1.0
13:29:19 WRN This is a warn log env=dev version=v0.1.0
13:29:19 ERR This is an error log env=dev version=v0.1.0
Welcome to k8s-controller-tutorial CLI!
```

and reading a variable from a file 'config.yaml'
```
log-level: error
```
$ ./controller
```
Using config file: /workspaces/k8s-controllers/config.yaml
13:34:21 ERR This is an error log env=dev version=v0.1.0
Welcome to k8s-controller-tutorial CLI!
```