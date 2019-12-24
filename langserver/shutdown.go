package langserver

import (
	"context"
	"os"

	lsp "github.com/sourcegraph/go-lsp"
)

func Shutdown(ctx context.Context, vs lsp.None) error {
	err := os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	return nil
}
