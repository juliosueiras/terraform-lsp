package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	lsp "github.com/sourcegraph/go-lsp"
)

func TextDocumentDidChange(ctx context.Context, vs lsp.DidChangeTextDocumentParams) error {
	tempFile.Truncate(0)
	tempFile.Seek(0, 0)
	tempFile.Write([]byte(vs.ContentChanges[0].Text))

	uri, err := absolutePath(string(vs.TextDocument.URI))
	if err != nil {
		return err
	}
	fileURL := uri.Filename()

	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(tempFile.Name(), fileURL)

	if !isTCP {
		TextDocumentPublishDiagnostics(StdioServer, ctx, lsp.PublishDiagnosticsParams{
			URI:         vs.TextDocument.URI,
			Diagnostics: DiagsFiles[fileURL],
		})
	}
	return nil
}
