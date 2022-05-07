.PHONY : all cli test

all: cli

cli:
	go build -o ./bin/demoChainCLI ./cmd/cli/

test:
	go test -v  ./test/