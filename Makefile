tidy:
	go mod tidy

build: tidy
	go build -o bin/cacheengine

run: build
	./bin/cacheengine
