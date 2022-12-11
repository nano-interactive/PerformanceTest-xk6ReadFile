MAKEFLAGS += --silent

all: clean build


## clean: Removes any previously created build artifacts.
clean:
	rm -f ./k6

## build: Builds a custom 'k6' with the local extension. 
build:
	go install go.k6.io/xk6/cmd/xk6@v0.7.0
	xk6 build --with github.com/radepopovic/xk6-read-file=.

## test: Executes any unit tests.
test:
	go test -cover -race ./...

.PHONY: build clean test
