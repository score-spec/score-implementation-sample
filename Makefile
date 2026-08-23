# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

## Display a list of the documented make targets
.PHONY: help
help:
	@echo Documented Make targets:
	@perl -e 'undef $$/; while (<>) { while ($$_ =~ /## (.*?)(?:\n# .*)*\n.PHONY:\s+(\S+).*/mg) { printf "\033[36m%-30s\033[0m %s\n", $$2, $$1 } }' $(MAKEFILE_LIST) | sort

.PHONY: .FORCE
.FORCE:

build:
	go build ./cmd/score-implementation-sample/

test:
	go vet ./...
	go env -w GOTOOLCHAIN="$(shell go env GOVERSION)+auto"
	go test ./... -cover -race

test-app: build
	./score-implementation-sample --version
	./score-implementation-sample init
	cat score.yaml
	./score-implementation-sample generate score.yaml
	cat manifests.yaml

build-container:
	docker build -t score-implementation-sample:local .

# The image runs as the nonroot user of the hardened base image, which does not own the
# mounted directory, so run the commands that write to it as the current user instead.
test-container: build-container
	docker run --rm score-implementation-sample:local --version
	docker run --rm --user $(shell id -u):$(shell id -g) -v .:/score-implementation-sample score-implementation-sample:local init
	cat score.yaml
	docker run --rm --user $(shell id -u):$(shell id -g) -v .:/score-implementation-sample score-implementation-sample:local generate score.yaml
	cat manifests.yaml