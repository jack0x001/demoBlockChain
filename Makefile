.PHONY : all cli test run

all: cli

cli:
	go build -o ./bin/demoChainCLI ./cmd/cli/

test:
	go test -v  ./test/

run:
	./bin/demoChainCLI run --ip=127.0.0.1 --port=51239 --datadir=./test/_bootstrap/

run1:
	./bin/demoChainCLI run --ip=127.0.0.1 --port=51240 --datadir=./test/_testdata/

run2:
	./bin/demoChainCLI run --ip=127.0.0.1 --port=51241 --datadir=./test/_testdata1/
