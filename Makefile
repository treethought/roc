help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SHELL := /bin/bash


plugins:
	go build -o ./bin/greeter ./examples/greeter
	go build -o ./bin/namer ./examples/namer
	go build -o ./dispatcher/dispatcher ./dispatcher/main.go


start:
	go run cmd/main.go run -c examples/config.yaml

run: plugins start
