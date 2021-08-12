.PHONY: all
all: build

build:
	go build -o ./bin/auto-launcher cmd/auto-launcher/main.go

test:
	go test ./... -v -count=1 -coverprofile=coverage.txt -covermode=atomic
