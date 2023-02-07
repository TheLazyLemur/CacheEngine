tidy:
	go mod tidy

test:
	go test -v ./...
	go test -v ./... --race

build: tidy test
	go build -o bin/cacheengine

run: build
	./bin/cacheengine
