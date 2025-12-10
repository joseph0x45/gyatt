build:
	go build .

test:
	rm -rf testing
	mkdir -p testing
	cd testing && ../gyatt init testing
