PACKAGE := github.com/hectcastro/heimdall

GPM_VERSION := v1.3.2
GPM_URL := https://raw.githubusercontent.com/pote/gpm/$(GPM_VERSION)/bin/gpm

GOX_OSARCH ?= linux/amd64 darwin/amd64
GOX_FLAGS ?= -output="pkg/{{.OS}}_{{.Arch}}/heimdall" -osarch="$(GOX_OSARCH)"

all: test

clean:
	rm -rf pkg/ gpm

test:
	@go test ./... -timeout=30s -parallel=4
	@go tool vet .

docker-test:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall sh -c 'sleep 1 && make test'

gpm:
	@wget -qN $(GPM_URL)
	@chmod +x $@

deps: gpm
	@./gpm install

gox-bootstrap:
	@gox -build-toolchain -osarch="$(GOX_OSARCH)"

release: deps gox-bootstrap
	@gox $(GOX_FLAGS) $(PACKAGE)

	@tar cvzf pkg/darwin_amd64/heimdall.tar.gz pkg/darwin_amd64/heimdall
	@tar cvzf pkg/linux_amd64/heimdall.tar.gz pkg/linux_amd64/heimdall

docker-release:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall gox $(GOX_FLAGS) $(PACKAGE)

.PHONY: all clean test docker-test gpm deps gox-bootstrap release docker-release
