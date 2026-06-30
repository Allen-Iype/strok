BINARY := strok
PKG := ./cmd/strok

.PHONY: build run test vet fmt clean

build:
	go build -o $(BINARY) $(PKG)

run:
	go run $(PKG)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
