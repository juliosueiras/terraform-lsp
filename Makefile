DATE    := $$(date +%Y-%m-%dT%T%z)
VERSION := $$(git describe --tags)
COMMIT  := $$(git rev-list -1 HEAD)
DST     ?= ~/.bin/

default:
	go build -ldflags "-X main.GitCommit=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE)"

copy:
	go build -ldflags "-X main.GitCommit=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE)" && cp ./terraform-lsp $(DST) && cp ./terraform-lsp ~/

.PHONY: copy
