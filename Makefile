help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SHELL := /bin/bash

build-cli:
	go build -o ./roc ./cmd/...

clean:
	rm -rf bin/*

build:
	go build -o ./bin/greeter ./examples/greeter
	go build -o ./bin/namer ./examples/namer
	go build -o ./bin/std/ ./std/...

protos:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/request.proto

start:
	go run cmd/main.go run -c examples/config.yaml

run: build start
