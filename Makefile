.PHONY: all
all: build

build:
	go build -o ./bin/auto-launcher cmd/al/main.go
