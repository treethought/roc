help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SHELL := /bin/bash


plugins:
	go build -o ./plugin/greeter ./plugin/greeter/greeter.go
	go build -o ./plugin/namer ./plugin/namer/namer.go
	go build -o ./dispatcher/dispatcher ./dispatcher/main.go


start:
	go run cmd/main.go

run: plugins start
