package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/memfs"
	log "github.com/sirupsen/logrus"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/spf13/afero"
)

func Initialize(ctx context.Context, vs lsp.InitializeParams) (lsp.InitializeResult, error) {
	file, err := afero.TempFile(memfs.MemFs, "", "tf-lsp-")
	if err != nil {
		log.Fatal(err)
	}
	//defer os.Remove(file.Name())
	tempFile = file

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    lsp.TDSKFull,
				},
			},
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{"."},
			},
			HoverProvider: false,
		},
	}, nil
}
