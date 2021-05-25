MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.4.2.2-DEV
CURARCH		:= $(shell uname -m | tr 'A-Z' 'a-z')
CUROS		:= $(shell uname -s | tr 'A-Z' 'a-z')
CURBRANCH	:= $(shell git branch | grep '*' | awk '{print $$2}')
BUILDARCH	:= $(CUROS)_$(CURARCH)
BUILDDIR	:= ./bin
TZ			:= UTC
BUILDVER    := -X github.com/wneessen/sotbot/version.Version=$(CURVER)
CURUSER     := $(shell whoami)
BUILDUSER   := -X github.com/wneessen/sotbot/version.BuildUser=$(subst $(SPACE),_,$(CURUSER))
CURDATE     := $(shell date +'%Y-%m-%d %H:%M:%S')
BUILDDATE   := -X github.com/wneessen/sotbot/version.BuildDate=$(subst $(SPACE),_,$(CURDATE))
BUILDBRANCH	:= -X github.com/wneessen/sotbot/version.BuildBranch=$(subst $(SPACE),_,$(CURBRANCH))
VEROS		:= -X github.com/wneessen/sotbot/version.BuildOs=$(subst $(SPACE),_,$(CUROS))
VERARCH		:= -X github.com/wneessen/sotbot/version.BuildArch=$(subst $(SPACE),_,$(CURARCH))
OUTFILE 	:= sotbot_$(CUROS)_$(CURARCH)_v$(CURVER)
TARGETS		:= clean build

all: $(TARGETS)

test:
	go test $(MODNAME)

clean: clean
release: clean build release clean

run:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(BUILDBRANCH) $(VERARCH) $(VEROS)" $(MODNAME)/cmd/sotbot

build:
	/usr/bin/env CGO_ENABLED=1 go build -o $(BUILDDIR)/v$(CURVER)/$(CUROS)/$(CURARCH)/sotbot -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(BUILDBRANCH) $(VERARCH) $(VEROS)" $(MODNAME)/cmd/sotbot
	ln -s v$(CURVER)/$(CUROS)/$(CURARCH)/sotbot $(BUILDDIR)/sotbot

clean:
	rm -rf bin/v$(CURVER)/$(CUROS)/$(CURARCH)
	rm -f $(BUILDDIR)/sotbot
	rm -f releases/v$(CURVER)/$(OUTFILE).tar.gz releases/v$(CURVER)/$(OUTFILE).tar.gz.sha256

release:
	cp -r LICENSE README.md documentation media config bin/v$(CURVER)/$(CUROS)/$(CURARCH)/
	mkdir -p releases/v$(CURVER)/
	tar czf releases/v$(CURVER)/$(OUTFILE).tar.gz bin/v$(CURVER)/$(CUROS)/$(CURARCH)/
	minisign -Sm releases/v$(CURVER)/$(OUTFILE).tar.gz -t "SoTBot v$(CURVER) - OS: $(CUROS) // Arch: $(CURARCH)"