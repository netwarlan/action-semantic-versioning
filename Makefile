BINARY_NAME = action-semantic-versioning

.PHONY: build test lint docker clean

build:
	go build -ldflags="-w -s" -o bin/$(BINARY_NAME) ./cmd/action-semantic-versioning

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

docker:
	docker build -t $(BINARY_NAME):latest .

clean:
	rm -rf bin/ coverage.out
