.PHONY: all

commands = auto-launcher auto-builder

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o ./bin/$@ $<

lint:
	golangci-lint run

test:
	go test ./... -v -count=1 -coverprofile=coverage.txt -covermode=atomic

cov-html:
	go tool cover -html=coverage.txt
