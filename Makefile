DATE    := $$(date +%Y-%m-%dT%T%z)
VERSION := $$(git describe --tags)
COMMIT  := $$(git rev-list -1 HEAD)
DST     ?= ~/.bin/

terraform-lsp:
	go build -ldflags "-X main.GitCommit=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE)"

copy: terraform-lsp | create-dir
	cp ./terraform-lsp $(DST) && cp ./terraform-lsp ~/

create-dir:
	mkdir -p $(DST)

clean:
	rm -f terraform-lsp

default: terraform-lsp

.PHONY: clean copy
