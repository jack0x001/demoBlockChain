.PHONY : all cli test run

all: cli

cli:
	go build -o ./bin/demoChainCLI ./cmd/cli/

test:
	go test -v  ./test/

run:
	./bin/demoChainCLI status --datadir=./database
	./bin/demoChainCLI balances list --datadir=./database