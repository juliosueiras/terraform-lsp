package langserver

import (
	"context"
	lsp "github.com/sourcegraph/go-lsp"
)

func TextDocumentDidClose(ctx context.Context, vs lsp.DidCloseTextDocumentParams) error {
	return nil
}
