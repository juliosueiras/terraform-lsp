default:
	go build -ldflags "-X main.GitCommit=$$(git rev-list -1 HEAD) -X main.Version=$$(git describe --tags) -X main.Date=$$(date --rfc-3339=date)" 

copy:
	go build -ldflags "-X main.GitCommit=$$(git rev-list -1 HEAD) -X main.Version=$$(git describe --tags) -X main.Date=$$(date --rfc-3339=date)" && cp ./terraform-lsp ~/.bin/ && cp ./terraform-lsp ~/

.PHONY: copy
