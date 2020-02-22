package langserver

import (
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
)

func CreateServer() *jrpc2.Server {
	Server = jrpc2.NewServer(handler.Map{
		"initialize":              handler.New(Initialize),
		"textDocument/completion": handler.New(TextDocumentComplete),
		"textDocument/didChange":  handler.New(TextDocumentDidChange),
		"textDocument/didOpen":    handler.New(TextDocumentDidOpen),
		"textDocument/didClose":   handler.New(TextDocumentDidClose),
		"textDocument/documentLink": handler.New(TextDocumentDocumentLink),
		//"textDocument/hover":      handler.New(TextDocumentHover),
		//"textDocument/references": handler.New(TextDocumentReferences),
		//"textDocument/codeLens": handler.New(TextDocumentCodeLens),
		"exit":            handler.New(Exit),
		"shutdown":        handler.New(Shutdown),
		"$/cancelRequest": handler.New(CancelRequest),
	}, &jrpc2.ServerOptions{
		AllowPush: true,
	})

	return Server
}
