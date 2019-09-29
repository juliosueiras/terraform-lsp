package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	lsp "github.com/sourcegraph/go-lsp"
	"strings"
)

func TextDocumentDidChange(ctx context.Context, vs lsp.DidChangeTextDocumentParams) error {
	tempFile.Truncate(0)
	tempFile.Seek(0, 0)
	tempFile.Write([]byte(vs.ContentChanges[0].Text))

	fileURL := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)
	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(tempFile.Name(), fileURL)

	TextDocumentPublishDiagnostics(Server, ctx, lsp.PublishDiagnosticsParams{
		URI:         vs.TextDocument.URI,
		Diagnostics: DiagsFiles[fileURL],
	})
	return nil
}
