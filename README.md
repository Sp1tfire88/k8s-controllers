# k8s-controllers
FastHTTP

## set dependencies
```
go get github.com/valyala/fasthttp
go get github.com/fasthttp/router
```
## build
```
go build -o controller
```

## run
$ ./controller server
```
Using config file: /workspaces/k8s-controllers/config.yaml
13:55:43 INF Starting FastHTTP server on :8080 env=dev version=v0.1.0
```



# Step4*: Add http requests logging
## set dependencies
```
github.com/google/uuid
```
## testing
```
$ go run main.go server --log-level=debug
```
```
curl  http://localhost:9000/health
curl -X POST http://localhost:9000/post -d '{"foo":"bar"}' -H "Content-Type: application/json"
```

## output
```
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-20T15:07:21Z","message":"Starting FastHTTP server on :9000"}
{"level":"info","env":"dev","version":"v0.1.0","method":"GET","path":"/health","remote_ip":"127.0.0.1","request_id":"da8b15ca-2531-46e2-ac0e-c1e2cee051c4","latency":0.009411,"time":"2025-06-20T15:07:54Z","message":"Request handled"}
{"level":"info","env":"dev","version":"v0.1.0","method":"GET","path":"/health","remote_ip":"127.0.0.1","request_id":"4aef3c35-1b89-4954-aea2-af9600c7ac7b","latency":0.005763,"time":"2025-06-20T15:08:51Z","message":"Request handled"}
{"level":"info","env":"dev","version":"v0.1.0","body":"{\"foo\":\"bar\"}","time":"2025-06-20T15:09:54Z","message":"Received POST data"}
{"level":"info","env":"dev","version":"v0.1.0","method":"POST","path":"/post","remote_ip":"127.0.0.1","request_id":"fcd4d4ae-0369-4f9f-9e48-6d06f6c2289b","latency":0.043429,"time":"2025-06-20T15:09:54Z","message":"Request handled"}
```
or set log-info and port in `config.yaml`
```
log-level: debug
port: 9000
```
```
$ go run main.go server
Using config file: /workspaces/k8s-controllers/config.yaml
{"level":"info","env":"dev","version":"v0.1.0","time":"2025-06-20T15:22:38Z","message":"Starting FastHTTP server on :9000"}
```