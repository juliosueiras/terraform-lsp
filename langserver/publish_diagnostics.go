package langserver

import (
	"context"

	"github.com/creachadair/jrpc2"
	lsp "github.com/sourcegraph/go-lsp"
)

func TextDocumentPublishDiagnostics(ctx context.Context, vs lsp.PublishDiagnosticsParams) error {

	var resultedError error

	if isTCP {
		resultedError = jrpc2.ServerFromContext(ctx).Notify(ctx, "textDocument/publishDiagnostics", vs)
	} else {
		resultedError = StdioServer.Notify(ctx, "textDocument/publishDiagnostics", vs)
	}

	return resultedError
}
