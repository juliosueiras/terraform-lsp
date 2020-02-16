package langserver

import (
	"context"
  "log"
  "github.com/juliosueiras/terraform-lsp/memfs"

	lsp "github.com/sourcegraph/go-lsp"
)

func Shutdown(ctx context.Context, vs lsp.None) error {
	err := memfs.MemFs.Remove(tempFile.Name())
	if err != nil {
		return err
	}

  log.Println("Shutdown")
	return nil
}
