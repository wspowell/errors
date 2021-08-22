CURRENT_DIR = $(shell pwd)

prereq:
	go install gotest.tools/gotestsum@latest

build:
	go build ./...
	go build -tags release ./...

test-debug: build
	gotestsum --format dots -- -count=1 -parallel 8 -race -cover -coverprofile=debug.cover -v ./...
	@go tool cover -func debug.cover | grep total | awk '{print "Coverage "substr($$3, 1, length($$3)-1)"%"}'

test-release: build
	gotestsum --format dots -- -count=1 -parallel 8 -race -cover -coverprofile=release.cover -v -tags release ./...
	@go tool cover -func release.cover | grep total | awk '{print "Coverage "substr($$3, 1, length($$3)-1)"%"}'

test: test-debug test-release

bench: build
	# Run benchmarks with -race for testing purposes (since -race adds overhead to real benchmarks).
	go test -bench=. -benchmem -count=1 -parallel 8 -race
	go test -bench=. -benchmem -count=1 -parallel 8 -tags release -race
	#
	# *** Run for real ***
	#
	go test -bench=. -benchmem -count=1 -parallel 8 
	go test -bench=. -benchmem -count=1 -parallel 8 -tags release