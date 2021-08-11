.PHONY: all
all: build

build:
	go build -o ./bin/auto-launcher cmd/auto-launcher/main.go
