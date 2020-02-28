package langserver

import (
	"context"
)

type InitializedParams struct{}

func Initialized(ctx context.Context, vs InitializedParams) error {
	return nil
}
