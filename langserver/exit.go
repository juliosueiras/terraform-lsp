package langserver

import (
	"context"
	lsp "github.com/sourcegraph/go-lsp"
	"os"
)

func Exit(ctx context.Context, vs lsp.None) error {
	os.Remove(tempFile.Name())
	return nil
}
