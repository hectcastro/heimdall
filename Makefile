PACKAGE = github.com/hectcastro/heimdall
PROJECT_PACKAGES = $(shell go list ./... | grep -v /vendor/)
GOX_OSARCH ?= linux/amd64
GOX_FLAGS ?= -output="pkg/{{.OS}}_{{.Arch}}/heimdall" -osarch="$(GOX_OSARCH)"

all: test install

install:
	$(GOPATH)/bin/godep go install $(PROJECT_PACKAGES)

test:
	golint -set_exit_status .
	golint -set_exit_status heimdall
	$(GOPATH)/bin/godep go vet $(PROJECT_PACKAGES)
	$(GOPATH)/bin/godep go test -v $(PROJECT_PACKAGES) -timeout=30s -parallel=4

release:
	$(GOPATH)/bin/godep restore
	$(GOPATH)/bin/gox $(GOX_FLAGS) $(PACKAGE)

	tar -C pkg/linux_amd64 -cvzf pkg/linux_amd64_heimdall.tar.gz heimdall

# Docker

docker-test:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall sh -c 'sleep 1 && make test'

docker-release:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall sh -c 'godep restore && make release'

.PHONY: all test release docker-test docker-release
