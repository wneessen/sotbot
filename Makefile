MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.3.1
BUILDDIR	:= ./bin
TZ			:= UTC
BUILDVER    := -X github.com/wneessen/sotbot/version.Version=$(CURVER)
CURUSER     := $(shell whoami)
BUILDUSER   := -X github.com/wneessen/sotbot/version.BuildUser=$(subst $(SPACE),_,$(CURUSER))
CURDATE     := $(shell date +'%Y-%m-%d %H:%M:%S')
BUILDDATE   := -X github.com/wneessen/sotbot/version.BuildDate=$(subst $(SPACE),_,$(CURDATE))

TARGETS			:= build

all: $(TARGETS)

test:
	go test $(MODNAME)

build:
	/usr/bin/env CGO_ENABLED=1 go build -o $(BUILDDIR)/v$(CURVER)/sotbot -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot
	rm $(BUILDDIR)/sotbot
	ln -s $(BUILDDIR)/v$(CURVER)/sotbot $(BUILDDIR)/sotbot

build-linux-amd64:
	/usr/bin/env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/v$(CURVER)/linux/amd64/sotbot -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot
	rm $(BUILDDIR)/sotbot
	ln -s $(BUILDDIR)/v$(CURVER)/linux/amd64/sotbot $(BUILDDIR)/sotbot

run:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER)" $(MODNAME)/cmd/sotbot