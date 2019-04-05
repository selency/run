all: deps build

deps:
	go get -t ./...

clean:
	rm -f datadump

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags "-X main.version=1.0.0"
