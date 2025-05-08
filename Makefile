.PHONY: build build-all clean

build:
	go build -o pushmonitor cmd/pushmonitor/main.go

build-all: build-macos build-linux

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o build/pushmonitor-darwin-amd64 cmd/pushmonitor/main.go
	GOOS=darwin GOARCH=arm64 go build -o build/pushmonitor-darwin-arm64 cmd/pushmonitor/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/pushmonitor-linux-amd64 cmd/pushmonitor/main.go
	GOOS=linux GOARCH=arm64 go build -o build/pushmonitor-linux-arm64 cmd/pushmonitor/main.go

clean:
	rm -rf build/
	rm -f pushmonitor