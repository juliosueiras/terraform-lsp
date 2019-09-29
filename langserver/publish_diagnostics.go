package langserver

import (
	"context"
	"github.com/creachadair/jrpc2"
	lsp "github.com/sourcegraph/go-lsp"
)

func TextDocumentPublishDiagnostics(server *jrpc2.Server, ctx context.Context, vs lsp.PublishDiagnosticsParams) error {

	return server.Push(ctx, "textDocument/publishDiagnostics", vs)
}
