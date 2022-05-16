.PHONY : all cli test run

all: cli

cli:
	go build -o ./bin/demoChainCLI ./cmd/cli/

test:
	go test -v  ./test/

run: cli
	./bin/demoChainCLI run --datadir=./test/_testdata/
