# gRPC User Management Service

Based on this series of [tutorials](https://www.youtube.com/watch?v=YudT0nHvkkE&list=PLrSqqHFS8XPYu-elDr1rjbfk0LMZkAA4X)

## How to run this example

1. run the grpc server

```sh
$ go run server/server.go
```

2. run the client

```sh
$ go run client/client.go
```

## How to create proto files

1. use the makefile

```sh
$ make generate
```