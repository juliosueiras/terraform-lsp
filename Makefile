default:
	go build -ldflags "-X main.GitCommit=$$(git rev-list -1 HEAD) -X main.Version=$$(git describe --tags) -X main.Date=$$(date +'%Y-%m-%d')"

copy:
	go build -ldflags "-X main.GitCommit=$$(git rev-list -1 HEAD) -X main.Version=$$(git describe --tags) -X main.Date=$$(date +'%Y-%m-%d')" && cp ./terraform-lsp ~/.bin/ && cp ./terraform-lsp ~/

.PHONY: copy
