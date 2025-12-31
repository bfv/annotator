# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=annotator
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_LINUX=$(BINARY_NAME)
BUILDDIR=build
VERSION=dev-latest

.PHONY: all build build-windows build-linux clean

all: build

build: build-windows build-linux

build-windows:
	set GOOS=windows&& set GOARCH=amd64&& go build -o $(BUILDDIR)/$(BINARY_WINDOWS) -ldflags "-X main.version=$(VERSION) -w -s" ./cmd/annotator

build-linux:
	set GOOS=linux&& set GOARCH=amd64&& go build -o $(BUILDDIR)/$(BINARY_LINUX) -ldflags "-X main.version=$(VERSION) -w -s" ./cmd/annotator

clean:
	if exist annotator.exe del annotator.exe
	if exist annotator del annotator
	if exist annotations.json del annotations.json
	if exist annotations.log del annotations.log
