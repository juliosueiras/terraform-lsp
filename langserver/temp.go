package langserver

import (
	"github.com/creachadair/jrpc2"
	lsp "github.com/sourcegraph/go-lsp"
	"os"
)

var tempFile *os.File
var DiagsFiles = make(map[string][]lsp.Diagnostic)
var Server *jrpc2.Server
