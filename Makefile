.PHONY: resources

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
