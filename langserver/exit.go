package langserver

import (
	"context"
  "log"
  "os"
  "github.com/juliosueiras/terraform-lsp/memfs"

	lsp "github.com/sourcegraph/go-lsp"
)

func Exit(ctx context.Context, vs lsp.None) error {
	err := memfs.MemFs.Remove(tempFile.Name())
	if err != nil {
		return err
	}

  log.Println("Exited")
	os.Exit(0)
	return nil
}
