MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.3.2
BUILDDIR	:= ./bin
TZ			:= UTC
BUILDVER    := -X github.com/wneessen/sotbot/version.Version=$(CURVER)
CURUSER     := $(shell whoami)
BUILDUSER   := -X github.com/wneessen/sotbot/version.BuildUser=$(subst $(SPACE),_,$(CURUSER))
CURDATE     := $(shell date +'%Y-%m-%d %H:%M:%S')
BUILDDATE   := -X github.com/wneessen/sotbot/version.BuildDate=$(subst $(SPACE),_,$(CURDATE))
CURARCH		:= $(shell uname -sm | tr 'A-Z' 'a-z' | tr ' ' '-')
OUTFILE 	:= sotbot_$(CURARCH)_v$(CURVER)
TARGETS		:= clean build

all: $(TARGETS)

test:
	go test $(MODNAME)

build: build-$(CURARCH)
clean: clean-$(CURARCH)
release: clean build release-$(CURARCH) clean

build-linux-x86_64:
	/usr/bin/env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/v$(CURVER)/linux/amd64/sotbot -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot
	ln -s v$(CURVER)/linux/amd64/sotbot $(BUILDDIR)/sotbot

run:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot

clean-linux-x86_64:
	rm -rf bin/v$(CURVER)
	rm -f $(BUILDDIR)/sotbot
	rm -rf releases/v$(CURVER)

release-linux-x86_64:
	cp -r LICENSE README.md documentation media config bin/v$(CURVER)/linux/amd64/
	mkdir -p releases/v$(CURVER)/
	tar czf releases/v$(CURVER)/$(OUTFILE).tar.gz bin/v$(CURVER)/linux/amd64/
	sha256sum releases/v$(CURVER)/$(OUTFILE).tar.gz > releases/v$(CURVER)/$(OUTFILE).tar.gz.sha256

