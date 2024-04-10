run:
	@go run cmd/cli/main.go
.PHONY: run

build:
	@go build -o bin/wharf cmd/cli/main.go
.PHONY: build
