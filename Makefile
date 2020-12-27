TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')

default: test

test: fmtcheck
	go list $(TEST) | xargs -t -n4 go test $(TESTARGS) -v -timeout=2m -parallel=4

cover:
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck'"

.NOTPARALLEL:

.PHONY: cover default fmt fmtcheck test testrace 
