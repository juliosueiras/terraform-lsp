package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"bitbucket.org/creachadair/jrpc2"
	"bitbucket.org/creachadair/jrpc2/channel"
	"bitbucket.org/creachadair/jrpc2/handler"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	//"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/terraform/configs"
	//"github.com/hashicorp/terraform/providers"
	//"github.com/minamijoyo/tfschema/tfschema"
	"github.com/juliosueiras/terraform-lsp/hclstructs"
	"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	"github.com/sourcegraph/go-lsp"
)

var tempFile *os.File

var Server *jrpc2.Server

func Initialize(ctx context.Context, vs lsp.InitializeParams) (lsp.InitializeResult, error) {

	file, err := ioutil.TempFile("", "tf-lsp-")
	if err != nil {
		log.Fatal(err)
	}
	//defer os.Remove(file.Name())
	tempFile = file

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Options: &lsp.TextDocumentSyncOptions{
					OpenClose: true,
					Change:    1,
				},
			},
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{"."},
			},
			//			HoverProvider:             true,
			//			DocumentSymbolProvider:    true,
			//ReferencesProvider: true,
			//			DefinitionProvider:        true,
			//			DocumentHighlightProvider: true,
			//			CodeActionProvider:        true,
			//			RenameProvider:            true,
		},
	}, nil
}

func TextDocumentComplete(ctx context.Context, vs lsp.CompletionParams) (lsp.CompletionList, error) {
	parser := configs.NewParser(nil)

	fileURL := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)

	fileDir := filepath.Dir(fileURL)
	res, _ := filepath.Glob(fileDir + "/*.tf")
	var file *configs.File
	var resultFiles []*configs.File

	for _, v := range res {
		if fileURL == v {
			continue
		}

		cFile, _ := parser.LoadConfigFile(v)

		resultFiles = append(resultFiles, cFile)
	}

	column := -1
	var diags hcl.Diagnostics
	var hclFile *hclsyntax.Body
	file, diags, column, hclFile = helper.CheckAndGetConfig(parser, tempFile, vs.Position.Line+1, vs.Position.Character)

	resultFiles = append(resultFiles, file)

	files, diags := configs.NewModule(resultFiles, nil)

	fileText, _ := ioutil.ReadFile(tempFile.Name())
	pos := helper.FindOffset(string(fileText), vs.Position.Line+1, column)

	var result []lsp.CompletionItem

	posHCL := hcl.Pos{
		Byte: pos,
	}

	if r, found, _ := tfstructs.GetTypeCompletion(result, fileDir, hclFile, posHCL); found {
		helper.DumpLog("Found Type Completion")
		return r, nil
	}

	config, origConfig, configType := tfstructs.GetConfig(file, posHCL)

	if diags != nil || config == nil {
		helper.DumpLog("With Error or No Config")
		helper.DumpLog(diags)

		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        tfstructs.GetTopLevelCompletion(),
		}, nil
	}

	expr := config.OutermostExprAtPos(posHCL)
	attr := config.AttributeAtPos(posHCL)
	blocks := config.BlocksAtPos(posHCL)

	if expr == nil && attr == nil && blocks == nil {
		if r, found, _ := tfstructs.GetAttributeCompletion(result, configType, origConfig, fileDir); found {
			return r, nil
		}
	}

	// Block is Block, not resources
	//test := config.BlocksAtPos(posHCL)
	//helper.DumpLog(test)
	if blocks != nil && attr == nil {
		if r, found, _ := tfstructs.GetNestingCompletion(blocks, result, configType, origConfig, fileDir); found {
			return r, nil
		}
	}

	if expr != nil {
		origType := reflect.TypeOf(expr)
		if origType != hclstructs.ObjectConsExpr() {
			variables := hclstructs.GetExprVariables(origType, expr, posHCL)

			if len(variables) != 0 {
				if variables[0].RootName() == "var" {
					vars := variables[0]

					result = helper.ParseVariables(vars[1:], files.Variables, result)
				}

				return lsp.CompletionList{
					IsIncomplete: false,
					Items:        result,
				}, nil
			}
		} else {
			if blocks == nil && attr != nil {
				if r, found, _ := tfstructs.GetNestingAttributeCompletion(attr, result, configType, origConfig, fileDir, posHCL); found {
					return r, nil
				}
			}
		}
	}

	//spew.Dump(config)

	return lsp.CompletionList{
		IsIncomplete: false,
		Items:        result,
	}, nil

}

func TextDocumentDidChange(ctx context.Context, vs lsp.DidChangeTextDocumentParams) error {
	tempFile.Truncate(0)
	tempFile.Seek(0, 0)
	tempFile.Write([]byte(vs.ContentChanges[0].Text))
	fileURL := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)
	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(tempFile.Name(), fileURL)

	TextDocumentPublishDiagnostics(Server, ctx, lsp.PublishDiagnosticsParams{
		URI:         vs.TextDocument.URI,
		Diagnostics: DiagsFiles[fileURL],
	})
	return nil
}

var DiagsFiles = make(map[string][]lsp.Diagnostic)

func TextDocumentDidOpen(ctx context.Context, vs lsp.DidOpenTextDocumentParams) error {
	fileURL := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)
	DiagsFiles[fileURL] = tfstructs.GetDiagnostics(fileURL, fileURL)

	TextDocumentPublishDiagnostics(Server, ctx, lsp.PublishDiagnosticsParams{
		URI:         vs.TextDocument.URI,
		Diagnostics: DiagsFiles[fileURL],
	})
	tempFile.Write([]byte(vs.TextDocument.Text))
	return nil
}

func Exit(ctx context.Context, vs lsp.None) error {
	os.Remove(tempFile.Name())
	return nil
}

func TextDocumentDidClose(ctx context.Context, vs lsp.DidCloseTextDocumentParams) error {
	return nil
}

func CancelRequest(ctx context.Context, vs lsp.CancelParams) error {
	return nil
}

//func TextDocumentCodeLens(ctx context.Context, vs lsp.CodeLensParams) ([]lsp.CodeLens, error) {
//	return []lsp.CodeLens{
//		lsp.CodeLens{
//			Range: lsp.Range{
//				Start: lsp.Position{
//					Line:      7,
//					Character: 1,
//				},
//				End: lsp.Position{
//					Line:      7,
//					Character: 1,
//				},
//			},
//			Command: lsp.Command{
//				Title:   "References",
//				Command: "test",
//			},
//		},
//	}, nil
//}

//func TextDocumentReferences(ctx context.Context, vs lsp.ReferenceParams) ([]lsp.Location, error) {
//	return []lsp.Location{
//		lsp.Location{
//			URI: vs.TextDocument.URI,
//			Range: lsp.Range{
//				Start: lsp.Position{
//					Line:      3,
//					Character: 1,
//				},
//				End: lsp.Position{
//					Line:      3,
//					Character: 3,
//				},
//			},
//		},
//	}, nil
//}

func TextDocumentPublishDiagnostics(server *jrpc2.Server, ctx context.Context, vs lsp.PublishDiagnosticsParams) error {

	return server.Push(ctx, "textDocument/publishDiagnostics", vs)
}

func main() {
	Server = jrpc2.NewServer(handler.Map{
		"initialize":              handler.New(Initialize),
		"textDocument/completion": handler.New(TextDocumentComplete),
		"textDocument/didChange":  handler.New(TextDocumentDidChange),
		"textDocument/didOpen":    handler.New(TextDocumentDidOpen),
		"textDocument/didClose":   handler.New(TextDocumentDidClose),
		//"textDocument/references": handler.New(TextDocumentReferences),
		//"textDocument/codeLens": handler.New(TextDocumentCodeLens),
		"exit":            handler.New(Exit),
		"$/cancelRequest": handler.New(CancelRequest),
	}, &jrpc2.ServerOptions{
		AllowPush: true,
	})

	f, err := os.OpenFile("tf-lsp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// Start the server on a channel comprising stdin/stdout.
	Server.Start(channel.Header("")(os.Stdin, os.Stdout))
	log.Print("Server started")

	// Wait for the server to exit, and report any errors.
	if err := Server.Wait(); err != nil {
		log.Printf("Server exited: %v", err)
	}

	//	local := server.NewLocal(handler.Map{
	//		"initialize": handler.New(Initialize),
	//	}, &server.LocalOptions{
	//		ServerOptions: &jrpc2.ServerOptions{
	//			Logger:  log.New(os.Stderr, "[jhttp.Bridge] ", log.LstdFlags|log.Lshortfile),
	//			Metrics: metrics.New(),
	//		},
	//	})
	//
	//	http.Handle("/", jhttp.NewClientBridge(local.Client))
	//	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8888), nil))
}
