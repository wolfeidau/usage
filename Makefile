# Program version
VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')

# Add the godep path to the GOPATH
GOPATH=$(shell godep path):$(shell echo $$GOPATH)

deps:
	godep get

clean: deps
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test -v ./...

.PHONY: build dist clean test
