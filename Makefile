.PHONY: build test lint vet fmt clean run

BINARY := ai-reads

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

check: fmt vet lint test

clean:
	rm -f $(BINARY)
