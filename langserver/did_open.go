package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	lsp "github.com/sourcegraph/go-lsp"
	"strings"
)

func TextDocumentDidOpen(ctx context.Context, vs lsp.DidOpenTextDocumentParams) error {
	fileURL := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)
	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(fileURL, fileURL)

	TextDocumentPublishDiagnostics(Server, ctx, lsp.PublishDiagnosticsParams{
		URI:         vs.TextDocument.URI,
		Diagnostics: DiagsFiles[fileURL],
	})
	tempFile.Write([]byte(vs.TextDocument.Text))
	return nil
}
