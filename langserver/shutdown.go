package langserver

import (
	"context"
  "github.com/juliosueiras/terraform-lsp/memfs"

  log "github.com/sirupsen/logrus"
	lsp "github.com/sourcegraph/go-lsp"
)

func Shutdown(ctx context.Context, vs lsp.None) error {
	err := memfs.MemFs.Remove(tempFile.Name())
	if err != nil {
		return err
	}

  log.Info("Shutdown")
	return nil
}
