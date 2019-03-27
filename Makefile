VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := C123R
REPOPATH ?= $(ORG)/$(OWNER)/kubectl-aks

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(shell mkdir -p ./release)

.PHONY: build
build: release/kubectl-aks-$(GOOS)-$(GOARCH)

release/kubectl-aks-%-$(GOARCH): 
	CGO_ENABLED=0 GOOS=$* GOARCH=$(GOARCH) go build \
	  -a -o $@ cmd/kubectl-aks.go

.PHONY: dep 
dep:
	dep ensure

.PHONY: clean
clean:
	rm -rf release/