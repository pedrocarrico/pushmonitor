.PHONY: build build-all clean

VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X github.com/pedrocarrico/pushmonitor/internal/version.Version=$(VERSION) \
          -X github.com/pedrocarrico/pushmonitor/internal/version.BuildTime=$(BUILD_TIME) \
          -X github.com/pedrocarrico/pushmonitor/internal/version.GitCommit=$(GIT_COMMIT)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/pushmonitor ./cmd/pushmonitor

build-all: clean build-macos build-linux

build-macos:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o build/pushmonitor-darwin-amd64 cmd/pushmonitor/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o build/pushmonitor-darwin-arm64 cmd/pushmonitor/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o build/pushmonitor-linux-amd64 cmd/pushmonitor/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o build/pushmonitor-linux-arm64 cmd/pushmonitor/main.go

clean:
	rm -rf build/
	rm -f bin/pushmonitor

coverage_report: test
	go tool cover -html=coverage.out

fmt:
	gofmt -s -w .

test:
	go test ./... -v -coverprofile coverage.out
