GPM_VERSION := "v1.3.2"

all: test

clean:
	rm -rf pkg/*

test:
	@go test ./... -timeout=30s -parallel=4
	@go tool vet .

deps: gpm
	@./gpm install

release:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall \
		gox -output "pkg/{{.OS}}_{{.Arch}}/heimdall" \
				-osarch="linux/amd64 darwin/amd64"

	@tar cvzf pkg/darwin_amd64/heimdall.tar.gz pkg/darwin_amd64/heimdall
	@tar cvzf pkg/linux_amd64/heimdall.tar.gz pkg/linux_amd64/heimdall

github-release: release
	@github-release release \
		--user hectcastro \
		--repo heimdall \
		--tag $(BUILDKITE_TAG)

	@github-release upload \
		--user hectcastro \
		--repo heimdall \
		--tag $(BUILDKITE_TAG) \
		--name linux_amd64_heimdall.tar.gz \
		--file pkg/linux_amd64/heimdall.tar.gz

	@github-release upload \
		--user hectcastro \
		--repo heimdall \
		--tag $(BUILDKITE_TAG) \
		--name darwin_amd64_heimdall.tar.gz \
		--file pkg/darwin_amd64/heimdall.tar.gz

gpm:
	@wget -q https://raw.githubusercontent.com/pote/gpm/$(GPM_VERSION)/bin/gpm
	@chmod +x gpm

ci:
	@docker-compose build heimdall
	@docker-compose run --rm heimdall make test

.PHONY: all test deps release clean ci github-release
