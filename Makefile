SHELL = /bin/bash
PROJECT = hans

dependencies:
	go get -v -t ./...

test:
	go test -v ./...

build: clean
	go build -o build/hans

clean:
	rm -rf build/
