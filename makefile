CURRENT_DIR = $(shell pwd)

build:
	go build ./...

test: build 
	GOBIN=$(CURRENT_DIR)/bin GO111MODULE=off go get gotest.tools/gotestsum
	$(CURRENT_DIR)/bin/gotestsum --format dots -- -count=1 -parallel 8 -race -v ./...
	$(CURRENT_DIR)/bin/gotestsum --format dots -- -count=1 -parallel 8 -race -v -tags release ./...

bench: build
	go test  -bench=. -benchmem -count=1 -parallel 8 -race
	go test -tags release -bench=. -benchmem -count=1 -parallel 8 -race
