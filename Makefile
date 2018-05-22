GOCMD=go
GOBUILD=$(GOCMD) build

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OS:=linux
endif
ifeq ($(UNAME_S),Darwin)
	OS:=darwin
endif
ifndef OS
$(error Unknown OS)
endif

all: clean build

build: agent worker

linux: OS=linux
linux: build

macos: OS=darwin
macos: build

agent:
	GOOS=$(OS) $(GOBUILD) -o ./agent cmd/agent/main.go

worker:
	GOOS=$(OS) $(GOBUILD) -o ./worker cmd/worker/main.go

clean:
	rm -f ./agent ./worker
