
all:
	go build

install:
	go build
	cp mdserver $(HOME)/bin

neat:
	gofmt -w *.go
