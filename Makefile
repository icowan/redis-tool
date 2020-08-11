APPNAME = redis-tool
BIN = $(GOPATH)/bin
GOCMD = /usr/local/go/bin/go
GOBUILD = $(GOCMD) build

all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-windows-amd64

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./$(APPNAME)-linux-amd64 -v

build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o ./$(APPNAME)-linux-arm64 -v

build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ./$(APPNAME)-darwin-amd64 -v

build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o ./$(APPNAME)-windows-amd64.exe -v
