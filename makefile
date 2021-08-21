CURRENT_DIR = $(shell pwd)

prereq:
	go install gotest.tools/gotestsum@latest

build:
	go build ./...

test: build 
	gotestsum --format dots -- -count=1 -parallel 8 -race -v ./...
	gotestsum --format dots -- -count=1 -parallel 8 -race -v -tags release ./...

bench: build
	# Run benchmarks with -race for testing purposes (since -race adds overhead to real benchmarks).
	go test -bench=. -benchmem -count=1 -parallel 8 -race
	go test -bench=. -benchmem -count=1 -parallel 8 -tags release -race
	#
	# *** Run for real ***
	#
	go test -bench=. -benchmem -count=1 -parallel 8 
	go test -bench=. -benchmem -count=1 -parallel 8 -tags release