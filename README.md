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
