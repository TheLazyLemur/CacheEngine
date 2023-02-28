#!/bin/bash

set -xe

function tidy() {
	go mod tidy
}

function test() {
	go test -v ./...
	go test -v ./... --race
}

function build() {
	go build -o bin/cacheengine
}

function run() {
	./bin/cacheengine
}

command=$1

if [ $command = "tidy" ]; then
	tidy
elif [ $command = "test" ]; then
	test
elif [ $command = "build" ]; then
	tidy
	test
	build
elif [ $command = "run" ]; then
	test
	build
	run
fi
