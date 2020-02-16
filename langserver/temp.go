package langserver

import (
	"github.com/creachadair/jrpc2"
	"github.com/spf13/afero"
	lsp "github.com/sourcegraph/go-lsp"
)

var tempFile afero.File
var DiagsFiles = make(map[string][]lsp.Diagnostic)
var Server *jrpc2.Server
