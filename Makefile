# Makefile for mdserver
	
PROG = mdserver
VERSION = 0.0.2
TARDIR = 	$(HOME)/Desktop/TarPit/
DATE = 	`date "+%Y-%m-%d.%H_%M_%S"`
DOCOUT = README-$(PROG)-godoc.md

all:
	godatetime > compileDate.go
	go build -v

# change cp to echo if you really don't want to install the program
install:
	godatetime > compileDate.go
	go build -v
	go tool vet .
	go tool vet -shadow .
	gofmt -w *.go
	cp $(PROG) $(HOME)/bin
#	go install

# note that godepgraph can be used to derive .travis.yml install: section
docs:
	godoc2md . > $(DOCOUT)
	godepgraph -md -p . >> $(DOCOUT)
	deadcode -md >> $(DOCOUT)
	cp README-$(PROG).md README.md
	cat $(DOCOUT) >> README.md

neat:
	go fmt ./...

dead:
	deadcode > problems.dead

index:
	cindex .

clean:
	go clean ./...
	rm -f *~ problems.dead count.out
	rm -f $(DOCOUT) README2.md

tar:
	echo $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar
	tar -ncvf $(TARDIR)$(PROG)_$(VERSION)_$(DATE).tar .

