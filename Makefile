tidy:
	go mod tidy

test:
	go test -v ./...
	go test -v ./... --race

build: tidy test
	go build -o bin/cacheengine

run: build
	./bin/cacheengine

runfollower: build
	./bin/cacheengine --listenaddr :4000 --leaderaddr :3000
