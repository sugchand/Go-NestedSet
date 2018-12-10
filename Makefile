GO ?= go
GOPATH := $(CURDIR)
GOSRCPATH := ./src/NestedSet
GOBINPATH := ./bin
GOOUTPUTBIN := $(GOBINPATH)/NestedSet
SHELL := /bin/bash
DEP := $(shell command -v dep  2> /dev/null)

export GOPATH
export GOSRCPATH
export
.PHONY : clean build debug

all: build

clean:
	rm -rf $(GOBINPATH)/*

debug: GCFLAGS = -gcflags=all="-N -l"

debug:
	@echo -e "\nCompiling NestedSet app with gdb symbols..."
	$(MAKE) build

build:
ifndef DEP
$(error "dep is not available please install it with 'apt install go-dep'")
endif
	#Make sure the 'dep status' shows right data
	#-@(cd $(GOSRCPATH);$(DEP) status 2> /dev/null)
	@echo -e "\n\tSet 'GOPATH' to '$(GOPATH)'"
	@echo -e "\tRun 'env DEPNOLOCK=1 dep ensure' in $(GOSRCPATH) to install missing third party packages\n"
	$(GO) build $(GCFLAGS) -o $(GOOUTPUTBIN) $(GOSRCPATH)
	@echo -e "\n\t**** RESULT : $$? : Build completed!!! ****\n\t**** Binary is at $$PWD/bin ****"

tests:
	@echo -e "\n\tSet 'GOPATH' to '$(GOPATH)'"
	@echo -e "\tRun 'env DEPNOLOCK=1 dep ensure' in $(GOSRCPATH) to install missing third party packages\n"
	$(GO) test -v $(GOSRCPATH)/...
