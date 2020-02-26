package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/memfs"
	log "github.com/sirupsen/logrus"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/spf13/afero"
)

type DocumentLinkOptions struct {
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

type ExtendedServerCapabilities struct {
	TextDocumentSync     *lsp.TextDocumentSyncOptionsOrKind `json:"textDocumentSync,omitempty"`
	CompletionProvider   *lsp.CompletionOptions             `json:"completionProvider,omitempty"`
	HoverProvider        bool                               `json:"hoverProvider,omitempty"`
	DocumentLinkProvider *DocumentLinkOptions               `json:"documentLinkProvider,omitempty"`
}

type ExtendedInitializeResult struct {
	Capabilities ExtendedServerCapabilities `json:"capabilities"`
}

func Initialize(ctx context.Context, vs lsp.InitializeParams) (ExtendedInitializeResult, error) {
	file, err := afero.TempFile(memfs.MemFs, "", "tf-lsp-")
	if err != nil {
		log.Fatal(err)
	}
	//defer os.Remove(file.Name())
	tempFile = file

	return ExtendedInitializeResult{
		Capabilities: ExtendedServerCapabilities{
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
			DocumentLinkProvider: &DocumentLinkOptions{
				ResolveProvider: false,
			},
		},
	}, nil
}
