BINARY := crawler

.PHONY: build run vet tidy clean

build:
	go build -o bin/$(BINARY) ./cmd/crawler

run:
	go run ./cmd/crawler

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/ out