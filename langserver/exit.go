package langserver

import (
	"context"
	"os"

	lsp "github.com/sourcegraph/go-lsp"
)

func Exit(ctx context.Context, vs lsp.None) error {
	err := os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
