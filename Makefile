all: test

clean:
	rm -f heimdall_linux_amd64*
	rm -f heimdall_darwin_amd64*

test:
	@go test ./... -timeout=30s -parallel=4
	@go tool vet .

deps: gpm
	@./gpm install

release:
	@gox -osarch="linux/amd64 darwin/amd64"

gpm:
	@wget -q https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm
	@chmod +x gpm

ci:
	@docker-compose run heimdall make test

.PHONY: all test deps release clean ci
