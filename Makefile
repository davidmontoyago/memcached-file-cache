# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -count=1
GOGET=$(GOCMD) get
GO111MODULE=on
GOOS?=darwin
GOARCH=amd64

.PHONY: memcached

all: test build

build:
	go mod vendor
	$(GOBUILD) ./

test:
	$(GOTEST) ./pkg/... -v

clean:
	$(GOCLEAN)

fmt:
	$(GOCMD) fmt ./...

memcached:
	docker run --name memcached -p 11211:11211 -d memcached:1.5 -m 1000
	make test-memcached

test-memcached:
	$(GOCMD) run ./memcached/test.go

teardown-memcached:
	docker rm -f memcached

run:
	$(GOCMD) run ./main.go