MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.0.0b
BUILDDIR	:= ./bin
TZ			:= UTC

ifeq ($(OS), Windows_NT)
	OUTFILE	:= $(BUILDDIR)/sotbot.exe
else
	OUTFILE	:= $(BUILDDIR)/sotbot
endif

TARGETS			:= build-prod
DOCKERTARGETS	:= build-prod dockerize

all: $(TARGETS)

test:
	go test $(MODNAME)

build-prod:
	/usr/bin/env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(OUTFILE) -ldflags="-s -w" $(MODNAME)/cmd/sotbot

run-prod:
	@go run -ldflags="-s -w" $(MODNAME)/cmd/sotbot

dockerize:
	@sudo docker build -t sotbot:v$(CURVER) .