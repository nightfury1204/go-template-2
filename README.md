# Go Boilerplate

## Start Rest 

## Build
```bash
$ ./build.sh
or
$ make build
```

## Build proto
```bash
$ ./build-proto.sh
or
$ make build-proto
```

## Application binary
```bash
$ go-boilerplate serve (for rest)
$ go-boilerplate serve-grpc (for grpc)
$ go-boilerplate serve-grpc-rest (for grpc rest proxy)
or
$ make run
```

## Container dev
```bash
$ docker-compose up --build
or
$ make serve
```

## GuideLine

* api folder contains rest code
* rpcs folder contains grpc code
* rpcrestproxy contains grpc rest proxy code

* infra contains drivers like db, messaging, cache etc
* repo folder contains database code
* model folder contains model
* service folder contains application service

### flow
> cmd -> api/rpcs/rpcrestproxy -> service -> repo, models, cache, messaging