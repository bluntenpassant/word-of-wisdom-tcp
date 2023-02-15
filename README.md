# Word of Wisdom PoW TCP

# Using

```sh
go run ./cmd/server/server.go
go run ./cmd/client/client.go
```
or with docker
```sh
docker build -o word-of-wisdom-tcp ./docker
docker run word-of-wisdom-tcp -p 8081:8081
```