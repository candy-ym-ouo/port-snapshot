APP := port-snapshot
VERSION ?= dev
DIST := dist

.PHONY: test vet build release clean

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(APP) .

release:
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(APP)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(APP)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(APP)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(APP)-windows-amd64.exe .

clean:
	rm -rf bin dist
