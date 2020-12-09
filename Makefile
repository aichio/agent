# Go parameters
# # build with version infos
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

BINARY_NAME=agent

version = v1.0
buildDate = $(shell date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)

ldflags="-w -X main.verSion=$(version) -X main.buildDate=${buildDate} -X main.gitCommit=${gitCommit}"

all: clean build
build:
	$(GOBUILD) -ldflags ${ldflags} -o $(BINARY_NAME) -v
debug:
	$(GOBUILD) -gcflags "-N -l" -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
