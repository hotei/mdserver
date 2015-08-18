# Makefile for mdserver

PROG = mdserver

all:
	go build


# note godatetime is optional - if used it goes before 'go build'
#     if godatetime is not used you need to comment out where it's used in *.go
install:
	go tool vet .
	go tool vet -shadow .
	godatetime > compileDate.go
	go build
	gofmt -w *.go
	godoc2md . > README-$(PROG)-pkg-godoc.md
	godepgraph -md -p . >> README-$(PROG)-pkg-godoc.md
	cp mdserver $(HOME)/bin

neat:
	gofmt -w *.go
