```sh
go get --tool github.com/bufbuild/buf/cmd/buf@v1.50.0

go tool buf -v
go tool buf config init

go tool buf dep update
```

> Linting using buff and golang
> Buf to detect breaking changes

```sh
go run cmd/server/main.go
go run cmd/client/main.go
```

---

> [**Code & Learn**](https://www.youtube.com/watch?v=KyLv9XEM0DM&t=34s&ab_channel=Code%26Learn)

## [INTERCEPTORS](https://www.youtube.com/watch?v=TZShbZVNbLc&ab_channel=Code%26Learn)

```sh
Server Side Unary Interceptor
Client Side Unary Interceptor
Server Side Stream Interceptor
Client Side Stream Interceptor
```

```go
srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
    func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
        start := time.Now()

        log.Printf("unary call made with: +%v", info)
        reponse, err := handler(ctx, req)
        log.Printf("time taken: %s", time.Since(start))
        return reponse, err

    },
    func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
        log.Println("second interceptor")
        return handler(ctx, req)
    }),
    grpc.StreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
        log.Println("server side streaming interceptor")
        return handler(srv, ss)
    }),
)
```

```go
conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithChainUnaryInterceptor(
        func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
            return invoker(ctx, method, req, reply, cc, opts...)
        }),
    grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
        return streamer(ctx, desc, cc, method, opts...)
    }))
```

## [ProtoValidate](https://www.youtube.com/watch?v=BJDfB6CKgKs&ab_channel=Code%26Learn)

> `go tool buf dep update`

## [Health Checking, Graceful Shutdown, Service Config](https://www.youtube.com/watch?v=jATLdVAuMdg&ab_channel=Code%26Learn)
