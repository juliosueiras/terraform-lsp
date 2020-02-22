package langserver

import (
	"context"
	lsp "github.com/sourcegraph/go-lsp"
	log "github.com/sirupsen/logrus"
)

type documentLinkParams struct {
  TextDocument lsp.TextDocumentItem `json:"textDocument"`
}

type DocumentLink struct {
  Range lsp.Range `json:"range"`

  /**
  * The uri this link points to. If missing a resolve request is sent later.
  */
  Target string `json:"target"`

  Tooltip string `json:"tooltip"`
}

func TextDocumentDocumentLink(ctx context.Context, vs documentLinkParams) ([]DocumentLink, error) {
  log.Info(vs)

  return []DocumentLink{
    DocumentLink{
      Range: lsp.Range{
        Start: lsp.Position{
          Line: 1,
          Character: 1,
        },
        End: lsp.Position{
          Line: 1,
          Character: 10,
        },
      },
      Target: "https://github.com",
      Tooltip: "https://github.com",
    },
  }, nil
}
