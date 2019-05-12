UNAME := $(shell uname)

ifeq ($(UNAME), Darwin)
PLUGIN_NAME=kubectl-aks-darwin
endif

ifeq ($(UNAME), Linux)
PLUGIN_NAME=kubectl-aks
endif

linux:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o kubectl-aks cmd/kubectl-aks.go

windows:
	GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o kubectl-aks-windows cmd/kubectl-aks.go

darwin:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o kubectl-aks-darwin cmd/kubectl-aks.go

all: linux windows darwin

package:
	zip kubectl-aks.zip kubectl-aks kubectl-aks-windows kubectl-aks-darwin Makefile
