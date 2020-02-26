package langserver

import (
	"context"
  "fmt"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

func RunStdioServer() {
	isTCP = false

	StdioServer = jrpc2.NewServer(handler.Map{
		"initialize":                handler.New(Initialize),
		"textDocument/completion":   handler.New(TextDocumentComplete),
		"textDocument/didChange":    handler.New(TextDocumentDidChange),
		"textDocument/didOpen":      handler.New(TextDocumentDidOpen),
		"textDocument/didClose":     handler.New(TextDocumentDidClose),
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

	StdioServer.Start(channel.Header("")(os.Stdin, os.Stdout))

	log.Info("Server started")

	// Wait for the server to exit, and report any errors.
	if err := StdioServer.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}

	log.Info("Server Finish")
}

func RunTCPServer(port int) {
	isTCP = true

	ServiceMap = handler.Map{
		"initialize":                handler.New(Initialize),
		"textDocument/completion":   handler.New(TextDocumentComplete),
		"textDocument/didChange":    handler.New(TextDocumentDidChange),
		"textDocument/didOpen":      handler.New(TextDocumentDidOpen),
		"textDocument/didClose":     handler.New(TextDocumentDidClose),
		"textDocument/documentLink": handler.New(TextDocumentDocumentLink),
		//"textDocument/hover":      handler.New(TextDocumentHover),
		//"textDocument/references": handler.New(TextDocumentReferences),
		//"textDocument/codeLens": handler.New(TextDocumentCodeLens),
		"exit":            handler.New(Exit),
		"shutdown":        handler.New(Shutdown),
		"$/cancelRequest": handler.New(CancelRequest),
	}

	// Start the server on a channel comprising stdin/stdout.

	lst, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	newChan := channel.Header("")

	ctx := context.Background()

	ctx, cancelFunc := context.WithCancel(ctx)

	go func() {
		if err := server.Loop(lst, ServiceMap, &server.LoopOptions{
			Framing: newChan,
			ServerOptions: &jrpc2.ServerOptions{
				AllowPush: true,
			},
		}); err != nil {
			log.Errorf("Loop: unexpected failure: %v", err)
			cancelFunc()
			return
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("Server Finish")
		err := lst.Close()
		if err != nil {
			log.Info("Server Finish")
			log.Info(ctx.Err())
			return
		}
	}
}
