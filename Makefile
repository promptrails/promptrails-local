.PHONY: build run test lint vet fmt clean docker-build docker-run

VERSION ?= dev
BINARY := promptrails-local
LDFLAGS := -s -w -X main.Version=$(VERSION)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test -v -race ./...

lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
	rm -rf dist/

docker-build:
	docker build -t promptrails/local:dev .

docker-run: docker-build
	docker run -p 8080:8080 promptrails/local:dev

tidy:
	go mod tidy

validate: vet lint test build
