package langserver

import (
	"context"
	"github.com/juliosueiras/terraform-lsp/memfs"
	"os"

	log "github.com/sirupsen/logrus"
	lsp "github.com/sourcegraph/go-lsp"
)

func Exit(ctx context.Context, vs lsp.None) error {
	err := memfs.MemFs.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	log.Info("Exited")
	os.Exit(0)
	return nil
}
