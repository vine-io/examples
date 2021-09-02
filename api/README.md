启动服务
```shell
go run main.go
```

测试
```shell
curl "http://127.0.0.1:8080/api/v1/echo/lack" -H  "accept: application/json" -H  "Content-Type: application/json" 
> {"reply":"reply: lack"}%
```