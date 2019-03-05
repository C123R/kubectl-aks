VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

ORG := github.com
OWNER := C123R
REPOPATH ?= $(ORG)/$(OWNER)/kubectl-aks

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(shell mkdir -p ./out)

.PHONY: build
build: out/kubectl-aks-$(GOOS)-$(GOARCH)

out/kubectl-aks-%-$(GOARCH): 
	CGO_ENABLED=0 GOOS=$* GOARCH=$(GOARCH) go build \
	  -ldflags="-s -w -X $(REPOPATH)/pkg/version.version=$(VERSION)" \
	  -a -o $@ .

.PHONY: dep 
dep:
	dep ensure

.PHONY: clean
clean:
	rm -rf out/