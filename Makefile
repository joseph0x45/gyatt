.PHONY: resources

VERSION := $(shell git describe --tags --always --dirty)

build:
	go build .

test:
	rm -rf testing
	mkdir -p testing
	cd testing && ../gyatt init testing

resources:
	tar -xf resources.tar.xz

compress-resources:
	tar -cJf resources.tar.xz resources/

release:
	GOOS=linux GOARCH=amd64 \
	go build -ldflags "-X main.version=$(VERSION)" -o gyatt

	tar -cJf gyatt-$(VERSION)-linux-amd64.tar.xz gyatt
	rm -f gyatt
