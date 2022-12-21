.PHONY: all
all: build

build:
	go build -o ./bin/auto-launcher ./cmd/auto-launcher/main.go

lint:
	golangci-lint run

test:
	go test ./... -v -count=1 -coverprofile=coverage.txt -covermode=atomic

cov-html:
	go tool cover -html=coverage.txt
