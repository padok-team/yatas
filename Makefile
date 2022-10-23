.DEFAULT_TARGET=help
VERSION:=$(shell cat VERSION)

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## build: Build Yatas binary
.PHONY: build
build: fmt vet
	go build -o bin/yatas

## fmt: Format source code
.PHONY: fmt
fmt:
	go fmt ./...

## vet: Vet source code
.PHONY: vet
vet:
	go vet ./...

## lint: Lint source code
.PHONY: lint
lint:
	golangci-lint run

## test: Run unit tests
.PHONY: test
test:
	go test ./... -cover

## Critic: Runs Go Critic
.PHONY: critic
critic:
	gocritic check ./...

## release: Release a new version
.PHONY: release
release: test
	standard-version
	git push --follow-tags origin main
