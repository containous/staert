.PHONY: clean dependencies checks test build fmt

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

default: clean checks test build

test: clean
	go test -v -race -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -f cover.out

build:
	GOOS=darwin go build;
	GOOS=windows go build;
	GOOS=linux go build;

checks:
	golangci-lint run

fmt:
	@gofmt -s -l -w $(SRCS)
