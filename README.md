#### Hexa RPC is Hexa-related RPC & gRPC SDK

#### Requirements:
go : minimum version `1.13`

#### Install
```
go get github.com/Kamva/hexa-rpc
```


#### Todo
- [ ] Use `recover` interceptor in of [gRPC interceptors](https://github.com/grpc-ecosystem/go-grpc-middleware).
- [ ] Implement status to Hexa error (and reverse) converter.
- [ ] Set Hexa logger as gRPC Logger (implement gRPC logger adapter by hexa logger)
- [ ] Collection presenter
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI