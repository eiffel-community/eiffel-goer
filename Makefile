# Install tools locally instead of in $HOME/go/bin.
export GOBIN := $(CURDIR)/bin
export PATH := $(GOBIN):$(PATH)

export RELEASE_VERSION ?= $(shell git describe --always)
export DOCKER_REGISTRY ?= registry.nordix.org/eiffel
export DEPLOY ?= goer

COMPILEDAEMON = $(GOBIN)/CompileDaemon
GIT = git
GOER = $(GOBIN)/goer
GOLANGCI_LINT = $(GOBIN)/golangci-lint
GOVVV = $(GOBIN)/govvv
MOCKGEN = $(GOBIN)/mockgen
PIGEON = $(GOBIN)/pigeon

GOLANGCI_LINT_VERSION = v1.43.0

.PHONY: all
all: test build start

.PHONY: gen-deps
gen-deps: $(MOCKGEN) $(PIGEON)

.PHONY: gen
gen: gen-deps
	go generate ./...

.PHONY: build
build: gen $(GOVVV)
	$(GOVVV) build -o $(GOER) ./cmd/goer

.PHONY: clean
clean:
	$(RM) $(GOER) $(GOVVV) $(MOCKGEN) $(PIGEON)
	docker-compose --project-directory . -f deploy/$(DEPLOY)/docker-compose.yml rm || true
	docker volume rm goer-volume || true

.PHONY: check
check: staticcheck test

.PHONY: staticcheck
staticcheck: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

.PHONY: test
test: gen
	go test -cover -timeout 30s -race $(shell go list ./... | grep -v test) 

# Start a development docker with a database that restarts on file changes.
.PHONY: start
start: $(COMPILEDAEMON) gen-deps
	docker-compose --project-directory . -f deploy/$(DEPLOY)/docker-compose.yml up

.PHONY: stop
stop:
	docker-compose --project-directory . -f deploy/$(DEPLOY)/docker-compose.yml down

# Build a docker using the production Dockerfile
.PHONY: docker
docker:
	docker build --build-arg revision=$(RELEASE_VERSION) -t $(DOCKER_REGISTRY)/$(DEPLOY):$(RELEASE_VERSION) -f ./deploy/$(DEPLOY)/Dockerfile .

.PHONY: push
push:
	docker push $(DOCKER_REGISTRY)/$(DEPLOY):$(RELEASE_VERSION)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: check-dirty
check-dirty:
	$(GIT) diff --exit-code HEAD

# Build dependencies

$(COMPILEDAEMON):
	mkdir -p $(dir $@)
	go install github.com/githubnemo/CompileDaemon@v1.3.0

$(GOLANGCI_LINT):
	mkdir -p $(dir $@)
	curl -sfL \
			https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)

$(GOVVV):
	mkdir -p $(dir $@)
	go install github.com/ahmetb/govvv@v0.3.0

$(MOCKGEN):
	mkdir -p $(dir $@)
	go install github.com/golang/mock/mockgen@v1.6.0

$(PIGEON):
	mkdir -p $(dir $@)
	go install github.com/mna/pigeon@v1.1.0
