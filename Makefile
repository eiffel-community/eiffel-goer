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

GOLANGCI_LINT_VERSION := v1.53.3
GOLANGCI_LINT_INSTALLATION_SHA256 := 060f1f3deb31b3d3b9515d691d9a776354cd63c7fcb5e036f18f0444cf2c934b
GOLANGCI_LINT_BINARY_SHA256 := 09237052ea9582630019182e890288ec155567fc949e6f329e3beb2c6e76b1b5

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

# Download the installation script for golangci-lint, verify its SHA-256 digest,
# run it if everything checks out, and verify the resulting binary.
$(GOLANGCI_LINT):
	mkdir -p $(dir $@)
	curl -sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		> $@.install-script-unverified
	echo "$(GOLANGCI_LINT_INSTALLATION_SHA256) $@.install-script-unverified" | sha256sum -c --quiet -
	sh -s -- -b $(dir $@) $(GOLANGCI_LINT_VERSION) < $@.install-script-unverified
	rm -f $@.install-script-unverified
	echo "$(GOLANGCI_LINT_BINARY_SHA256) $@" | sha256sum -c --quiet - || ( rm $@ ; exit 1 )

$(GOVVV):
	mkdir -p $(dir $@)
	go install github.com/ahmetb/govvv@v0.3.0

$(MOCKGEN):
	mkdir -p $(dir $@)
	go install github.com/golang/mock/mockgen@v1.6.0

$(PIGEON):
	mkdir -p $(dir $@)
	go install github.com/mna/pigeon@v1.1.0
