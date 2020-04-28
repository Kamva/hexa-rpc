#### Hexa RPC is Hexa-related RPC & gRPC SDK

#### Requirements:
go : minimum version `1.13`

#### Install
```
go get github.com/Kamva/hexa-rpc
```


#### Todo
- [ ] Use `recover` interceptor in the [gRPC interceptors](https://github.com/grpc-ecosystem/go-grpc-middleware).
- [ ] Implement status to Hexa error (and reverse) mapper.
- [x] Set Hexa logger as gRPC Logger (implement gRPC logger adapter by hexa logger)
- [ ] Implement request logger (log start-time, end-time, method, error,...)
- [ ] We should implement all of our interceptors for the Stream request/responses also (for now
 we just support Unary Request/responses).
- [ ] Collection presenter
- [ ] Write Tests
- [ ] Add badges to readme.
- [ ] CI
