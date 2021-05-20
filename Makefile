MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.2.1
BUILDDIR	:= ./bin
TZ			:= UTC
BUILDVER    := -X github.com/wneessen/sotbot/version.Version=$(CURVER)
CURUSER     := $(shell whoami)
BUILDUSER   := -X github.com/wneessen/sotbot/version.BuildUser=$(subst $(SPACE),_,$(CURUSER))
CURDATE     := $(shell date +'%Y-%m-%d %H:%M:%S')
BUILDDATE   := -X github.com/wneessen/sotbot/version.BuildDate=$(subst $(SPACE),_,$(CURDATE))

ifeq ($(OS), Windows_NT)
	OUTFILE	:= $(BUILDDIR)/sotbot.exe
else
	OUTFILE	:= $(BUILDDIR)/sotbot
endif

TARGETS			:= build-local
DOCKERTARGETS	:= build-docker dockerize docker-publish

all: $(TARGETS)

docker: $(DOCKERTARGETS)

test:
	go test $(MODNAME)

build-local:
	/usr/bin/env CGO_ENABLED=1 go build -o $(OUTFILE) -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot

run:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot

build-docker:
	/usr/bin/env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(OUTFILE) -ldflags="-s -w" $(MODNAME)/cmd/sotbot

run-prod:
	@go run -ldflags="-s -w" $(MODNAME)/cmd/sotbot

dockerize:
	@sudo docker build -t sotbot:v$(CURVER) .

docker-publish:
	@sudo docker tag sotbot:v$(CURVER) wneessen/sotbot:latest
	@sudo docker push wneessen/sotbot:latest