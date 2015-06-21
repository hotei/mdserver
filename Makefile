
all:
	go build
	go tool vet .
	go tool vet -shadow .


# note godatetime is optional
install:
	godatetime > compileDate.go
	go build
	gofmt -w *.go
	cp README-mdserver.md README.md
	cp mdserver $(HOME)/bin

neat:
	gofmt -w *.go
