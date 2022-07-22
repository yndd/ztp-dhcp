# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

VERSION ?= latest
REPO ?= ghcr.io/steiler
# IMAGE_TAG_BASE defines the docker.io namespace and part of the image name for remote images.
# This variable is used to construct full image tags for ndd packages.
IMAGE_TAG_BASE ?= $(REPO)/ztp-dhcp

# Package
PKG ?= $(IMAGE_TAG_BASE)-package


KUBECTL_NDD_VERSION ?= v0.2.20


help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


docker-build: ## Build docker image with the manager.
	DOCKER_BUILDKIT=1 docker build -t $(IMAGE_TAG_BASE) .

docker-push: docker-build ## Push docker image with the manager.
	docker push $(IMAGE_TAG_BASE)

.PHONY: package-build
package-build: kubectl-ndd ## build ndd package.
	rm -rf package/*.nddpkg
	cd package;PATH=$$PATH:$(LOCALBIN) kubectl ndd package build -t provider;cd ..

.PHONY: package-push
package-push: kubectl-ndd ## build ndd package.
	cd package;ls;PATH=$$PATH:$(LOCALBIN) kubectl ndd package push ${PKG};cd ..

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL_NDD ?= $(LOCALBIN)/kubectl-ndd


.PHONY: kubectl-ndd
kubectl-ndd: $(KUBECTL_NDD) ## Download kubectl-ndd locally if necessary.
$(KUBECTL_NDD): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/yndd/ndd-core/cmd/kubectl-ndd@$(KUBECTL_NDD_VERSION)  ;\


.PHONY: update-yndd-dependencies
update-yndd-dependencies:
	go get -d -u github.com/yndd/ztp-webserver@master


MOCKDIR = pkg/mocks

.PHONY: mocks-gen
mocks-gen: mocks-rm ## Generate mocks for all the defined interfaces.
	go install github.com/golang/mock/mockgen@latest
	mockgen -package=mocks -source=pkg/devices/device.go -destination=$(MOCKDIR)/device.go
	mockgen -package=mocks -source=pkg/backend/backend.go -destination=$(MOCKDIR)/backend.go
	mockgen -package=mocks -source=pkg/devices/devicemanagerhandler.go -destination=$(MOCKDIR)/devicemanagerhandler.go
	mockgen -package=mocks -source=pkg/devices/devicemanagerregistrator.go -destination=$(MOCKDIR)/devicemanagerregistrator.go
	mockgen -package=mocks -destination=$(MOCKDIR)/packetconn.go net PacketConn

.PHONY: mocks-rm
mocks-rm: ## remove generated mocks
	rm -rf $(MOCKDIR)/*

.PHONY: test
test: ## Run test with coverage
	go test -v -coverprofile coverage.out ./... -coverpkg=./...
	grep -v '/mocks/' coverage.out > coverage.tmp && mv coverage.tmp coverage.out
	go tool cover -html coverage.out -o coverage.html
