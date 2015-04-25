PACKAGE = github.com/hectcastro/heimdall
VERSION = '$(shell git describe --tags --always --dirty)'
GOVERSION = '$(shell go version)'
BUILDTIME = '$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'
LDFLAGS = -X main.version $(VERSION) -X main.goVersion $(GOVERSION) -X main.buildTime $(BUILDTIME)

GOX_OSARCH ?= linux/amd64 darwin/amd64
GOX_FLAGS ?= -output="pkg/{{.OS}}_{{.Arch}}/heimdall" -osarch="$(GOX_OSARCH)"

all: test install

install: godep
	${GOPATH}/bin/godep go install -ldflags "$(LDFLAGS)" ./...

clean:
	rm -rf pkg/

test: godep
	${GOPATH}/bin/godep go test -v ./... -timeout=30s -parallel=4

vendor: godep
	rm -rf Godep
	${GOPATH}/bin/godep save ./...

release: godep gox-bootstrap
	${GOPATH}/bin/godep restore
	${GOPATH}/bin/gox $(GOX_FLAGS) -ldflags "$(LDFLAGS)" $(PACKAGE)

	tar cvzf pkg/darwin_amd64/heimdall.tar.gz pkg/darwin_amd64/heimdall
	tar cvzf pkg/linux_amd64/heimdall.tar.gz pkg/linux_amd64/heimdall


# Gox

gox: ${GOPATH}/bin/gox

${GOPATH}/bin/gox:
	go get -u github.com/mitchellh/gox
	go install github.com/mitchellh/gox

gox-bootstrap: gox
	${GOPATH}/bin/gox -build-toolchain -osarch="$(GOX_OSARCH)"


# Godep

godep: ${GOPATH}/bin/godep

${GOPATH}/bin/godep:
	go get -u github.com/tools/godep
	go install github.com/tools/godep


# Docker

docker-test:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall sh -c 'sleep 1 && make test'

docker-release:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall sh -c 'godep restore && make gox && gox $(GOX_FLAGS) $(PACKAGE)'

.PHONY: all clean test godep vendor gox gox-bootstrap release docker-test docker-release
