CURRENT_DIR = $(shell pwd)

prereq:
	go install gotest.tools/gotestsum@latest

build:
	golangci-lint run

	go build ./...

test: build
	gotestsum --format dots -- -count=1 -parallel 8 -race -cover -coverprofile=debug.cover -v ./...
	@go tool cover -func debug.cover | grep total | awk '{print "Coverage "substr($$3, 1, length($$3)-1)"%"}'

	# Run benchmarks with -race for testing purposes (since -race adds overhead to real benchmarks).
	go test -bench=. -benchmem -count=1 -parallel 8 -race ./...

bench: build
	go test -bench=. -benchmem -count=1 -parallel 8 ./...