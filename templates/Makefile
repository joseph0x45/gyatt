build:
	templ generate && go generate && go build .

build-release:
	templ generate && go generate && go build -ldflags="-w -s"

clean:
	go clean
