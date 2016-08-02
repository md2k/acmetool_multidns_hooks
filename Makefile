EXEC_NAME = dns.hook
GO = GODEBUG=sbrk=1 GO15VENDOREXPERIMENT=1 go
GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --tags)"

build:
	@echo "Building Current OS Version"
	$(GO) build $(GOFLAGS) -o $(EXEC_NAME)

osx:
	@echo "Building OSX x64 Version"
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(EXEC_NAME)

linux:
	@echo "Building Linux x64 Version"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(EXEC_NAME)

win:
	@echo "Building Windows x64 Version"
	GOOS="windows" GOARCH=amd64 $(GO) build $(GOFLAGS) -o "$(EXEC_NAME).exe"

all:
	@echo "Building All versions As release archiver per platform:"
	GOOS="darwin" GOARCH=amd64 $(GO) build $(GOFLAGS) -o "$(EXEC_NAME)" ; tar cjf "$(EXEC_NAME)_macos_amd64.tar.bz2" $(EXEC_NAME) ; rm -f "$(EXEC_NAME)"
	GOOS="linux" GOARCH=amd64 $(GO) build $(GOFLAGS) -o "$(EXEC_NAME)" ; tar cjf "$(EXEC_NAME)_linux_amd64.tar.bz2" $(EXEC_NAME) ; rm -f "$(EXEC_NAME)"
	GOOS="windows" GOARCH=amd64 $(GO) build $(GOFLAGS) -o "$(EXEC_NAME).exe" ; tar cjf "$(EXEC_NAME)_win_amd64.tar.bz2" "$(EXEC_NAME).exe" ; rm -f "$(EXEC_NAME).exe"

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs
