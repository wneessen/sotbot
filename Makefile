MODNAME		:= github.com/wneessen/sotbot
SPACE		:= $(null) $(null)
CURVER		:= 1.4.9-DEV
CURARCH		:= $(shell uname -m | tr 'A-Z' 'a-z')
CUROS		:= $(shell uname -s | tr 'A-Z' 'a-z')
BUILDARCH	:= $(CUROS)_$(CURARCH)
BUILDDIR	:= ./bin
TZ			:= UTC
BUILDVER    := -X github.com/wneessen/sotbot/version.Version=$(CURVER)
CURUSER     := $(shell whoami)
BUILDUSER   := -X github.com/wneessen/sotbot/version.BuildUser=$(subst $(SPACE),_,$(CURUSER))
CURDATE     := $(shell date +'%Y-%m-%d %H:%M:%S')
BUILDDATE   := -X github.com/wneessen/sotbot/version.BuildDate=$(subst $(SPACE),_,$(CURDATE))
VEROS		:= -X github.com/wneessen/sotbot/version.BuildOs=$(subst $(SPACE),_,$(CUROS))
VERARCH		:= -X github.com/wneessen/sotbot/version.BuildArch=$(subst $(SPACE),_,$(CURARCH))
DEVGUILD	:= 843575000987336755
GUILDID		:= -X github.com/wneessen/sotbot/bot.GuildID=$(DEVGUILD)
OUTFILE 	:= sotbot_$(CUROS)_$(CURARCH)_v$(CURVER)
TARGETS		:= clean build

all: $(TARGETS)

test:
	go test $(MODNAME)

release: clean build release-pkg clean

dev:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(VERARCH) $(VEROS) $(GUILDID)" $(MODNAME)/cmd/sotbot

reset-dev:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(VERARCH) $(VEROS) $(GUILDID)" $(MODNAME)/cmd/sotbot -r

reset-global:
	/usr/bin/env CGO_ENABLED=1 go run -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(VERARCH) $(VEROS)" $(MODNAME)/cmd/sotbot -r


build:
	/usr/bin/env CGO_ENABLED=1 go build -o $(BUILDDIR)/v$(CURVER)/$(CUROS)/$(CURARCH)/sotbot -ldflags="-s -w $(BUILDVER) $(BUILDDATE) $(BUILDUSER) $(VERARCH) $(VEROS)" $(MODNAME)/cmd/sotbot
	ln -s v$(CURVER)/$(CUROS)/$(CURARCH)/sotbot $(BUILDDIR)/sotbot

clean:
	rm -rf bin/v$(CURVER)/$(CUROS)/$(CURARCH)
	rm -f $(BUILDDIR)/sotbot
	rm -f releases/v$(CURVER)/$(OUTFILE).tar.gz releases/v$(CURVER)/$(OUTFILE).tar.gz.sha256

release-pkg:
	cp -r LICENSE README.md documentation media config bin/v$(CURVER)/$(CUROS)/$(CURARCH)/
	mkdir -p releases/v$(CURVER)/
	tar czf releases/v$(CURVER)/$(OUTFILE).tar.gz bin/v$(CURVER)/$(CUROS)/$(CURARCH)/
	minisign -Sm releases/v$(CURVER)/$(OUTFILE).tar.gz -t "SoTBot v$(CURVER) - OS: $(CUROS) // Arch: $(CURARCH)"
