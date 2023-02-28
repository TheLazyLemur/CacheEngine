# Distributed Cache Golang

This is my attempt at building a distibuted cache and protocol in Golang with minimal dependencies.

## Dependencies

- [The Go Programming Language](https://go.dev/doc/install)
- Make or Bash
- You can also use the built in go commands the Makefile and bash script are for convenience

## Get Dependencies

- `build.sh tidy` or
- `make tidy`

## How To Run

- `build.sh run` or
- `make run`

## Build A Binary

- `build.sh build` or
- `make build`

## Run Tests

- `build.sh test` or
- `make test`

## Checklist of project

- [x] Implement cache
- [x] Implement protocol
- [x] Implement tcp server
- [x] Implement api server
- [x] Implement tcp client
- [ ] Implement message distribution to followers
- [ ] Implement Raft consensus
- [ ] Implement Database for cache backup
