package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	lsp "github.com/sourcegraph/go-lsp"
)

func TextDocumentDidOpen(ctx context.Context, vs lsp.DidOpenTextDocumentParams) error {
	uri, err := absolutePath(string(vs.TextDocument.URI))
	if err != nil {
		return err
	}
	fileURL := uri.Filename()

	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(fileURL, fileURL)

	if !isTCP {
		TextDocumentPublishDiagnostics(StdioServer, ctx, lsp.PublishDiagnosticsParams{
			URI:         vs.TextDocument.URI,
			Diagnostics: DiagsFiles[fileURL],
		})
	}
	tempFile.Write([]byte(vs.TextDocument.Text))
	return nil
}
