BINARY=spider

build:
	go build -o $(BINARY) main.go

run:
	go run main.go $(arg)

test:
	go test -v ./...