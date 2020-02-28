package langserver

import (
	"context"
	"fmt"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
	log "github.com/sirupsen/logrus"
	oldLog "log"
	"net"
	"os"
)

func RunStdioServer(oldLogInstance *oldLog.Logger) {
	isTCP = false

	StdioServer = jrpc2.NewServer(ServiceMap, &jrpc2.ServerOptions{
		AllowPush: true,
		Logger:    oldLogInstance,
	})

	StdioServer.Start(channel.Header("")(os.Stdin, os.Stdout))

	log.Info("Server started")

	// Wait for the server to exit, and report any errors.
	if err := StdioServer.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}

	log.Info("Server Finish")
}

func RunTCPServer(address string, port int, oldLogInstance *oldLog.Logger) {
	isTCP = true

	lst, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
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
				Logger:    oldLogInstance,
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

func InitializeServiceMap() {
	ServiceMap = handler.Map{
		"initialize":                handler.New(Initialize),
		"initialized":               handler.New(Initialized),
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
}
