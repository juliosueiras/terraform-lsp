package langserver

import (
	"github.com/creachadair/jrpc2"
	lsp "github.com/sourcegraph/go-lsp"
	"github.com/spf13/afero"
)

var tempFile afero.File
var DiagsFiles = make(map[string][]lsp.Diagnostic)
var StdioServer *jrpc2.Server
var ServiceMap jrpc2.Assigner
var isTCP bool
