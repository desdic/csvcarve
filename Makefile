.PHONY: all test

GOC=go
GOBIN=$(shell pwd)

all: test build

clean:
	rm -rf ./bin

build: bin/csvcarve

test:
	$(GOC) test -race -v ./...

bin/csvcarve: cmd/csvcarve/*.go
	CGO_ENABLED=0 GOBIN=${GOBIN}/bin $(GOC) install -mod=vendor -buildmode=pie ./cmd/csvcarve
