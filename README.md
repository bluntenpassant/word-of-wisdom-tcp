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

## Environment variables

### Client

| name           | type    | default        | description
|----------------|---------|----------------|--------------------------------------
| SERVER_ADDR    | string  | 127.0.0.1:9000 | listen TCP address
| FETCH_WORKERS  | int     | 4              | count of client requests at same time
| TIMEOUT        | int     | 1000           | timeout after failed request

### Server

| name             | type    | default        | description
|------------------|---------|----------------|----------------------------------------
| LISTEN_ADDR      | string  | 0.0.0.0:9000   | server TCP address
| DIFFICULTY       | byte    | 23             | difficulty of calc algorithm for client
| PROOF_TOKEN_SIZE | int     | 64             | data size for proof calc for client

# Implementation

## Project structure

- deploy - dockerfiles for build
- cmd
    - server - server side app
    - client - client side app
- internal
    - pow - implementation "Proof Of Work" algorithm based on sha256
    - client - implementation "Proof Of Work" client requests
    - server - implementation "Proof Of Work" server listener

## Challenge-Response protocol

- Client connected to server
- Server write connection to log
- Server send puzzle packet
  | offset             | name         | length
  | -------------------|--------------|---------------
  |                  0 | difficulty   | 1 byte
  |                  1 | token size   | 2 bytes
  |                  3 | rand token   | ProofTokenSize
- Client calculate proof based on difficulty and data (rand token)
- Client send proof nonce to server
- Server check proof based on difficulty, data (rand token) and received nonce
- If proof is valid
    - write to log
    - server send to client quote from “Word of Wisdom”
    - client print response from server
    - connection close
- If proof is not valid
    - write log log
    - connection close

## Algorithm

I chose the sha256 algorithm because:
- It present in the standard go library
- "nonce" number not too big
- Easy calculate zeroes in hash
- The difficulty of the calculation is enough to guard from ddos

### Calculation hash basis

| offset | name              | length
|--------|-------------------|----------------------------------
| 0      | nonce             | 8 bytes
| 8      | data (rand token) | ProofTokenSize (usually 64 bytes)

Benchmarks:

**cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz**

### Benchmark - Calculate Proof (used on client side)

| difficulty | average time     | alloc bytes   | alloc count |
| -----------|------------------|---------------|-------------|
| 5          | 724.9 ns/op	    | 112 B/op	    | 2 allocs/op |
| 10         | 559328 ns/op	    | 112 B/op	    | 2 allocs/op |
| 15         | 7799120 ns/op    | 112 B/op	    | 2 allocs/op |
| 20         | 196695777 ns/op  | 112 B/op	    | 2 allocs/op |
| 22         | 1129589102 ns/op | 128 B/op      | 4 allocs/op |
| 25         | 5314950709 ns/op	| 120 B/op	    | 3 allocs/op |

### Benchmark - Check Buf Proof (used on server side)

| difficulty | average time | alloc bytes | alloc count |
| -----------|--------------|-------------|-------------|
| 5          | 331.8 ns/op	| 0 B/op	    | 0 allocs/op |
| 10         | 336.6 ns/op	| 0 B/op	    | 0 allocs/op |
| 15         | 334.2 ns/op	| 0 B/op	    | 0 allocs/op |
| 20         | 331.4 ns/op	| 0 B/op	    | 0 allocs/op |
| 22         | 331.6 ns/op	| 0 B/op	    | 0 allocs/op |
| 25         | 332.2 ns/op  | 0 B/op	    | 0 allocs/op |


### Benchmark - Check Proof (with allocations)

| difficulty | average time | alloc bytes | alloc count |
| -----------|--------------|-------------|-------------|
| 5          | 370.6 ns/op  | 80 B/op     | 1 allocs/op |
| 10         | 379.3 ns/op  | 80 B/op     | 1 allocs/op |
| 15         | 372.8 ns/op  | 80 B/op     | 1 allocs/op |
| 20         | 368.4 ns/op  | 80 B/op     | 1 allocs/op |
| 22         | 375.2 ns/op  | 80 B/op     | 1 allocs/op |
| 25         | 368.9 ns/op  | 80 B/op     | 1 allocs/op |
