package langserver

import (
	"context"
	lsp "github.com/sourcegraph/go-lsp"
	"io/ioutil"
	"log"
)

func Initialize(ctx context.Context, vs lsp.InitializeParams) (lsp.InitializeResult, error) {

	file, err := ioutil.TempFile("", "tf-lsp-")
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
