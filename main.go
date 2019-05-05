package main

import (
	"context"
	//"fmt"
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
	"github.com/sourcegraph/go-lsp"
)

var tempFile *os.File

func Initialize(ctx context.Context, vs lsp.InitializeParams) (lsp.InitializeResult, error) {

	file, err := ioutil.TempFile("/tmp", "prefix")
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
			//			ReferencesProvider:        true,
			//			DefinitionProvider:        true,
			//			DocumentHighlightProvider: true,
			//			CodeActionProvider:        true,
			//			RenameProvider:            true,
		},
	}, nil
}

func TextDocumentComplete(ctx context.Context, vs lsp.CompletionParams) (lsp.CompletionList, error) {
	parser := configs.NewParser(nil)

	fileUrl := strings.Replace(string(vs.TextDocument.URI), "file://", "", 1)
	file2, _ := parser.LoadConfigFile(filepath.Dir(fileUrl) + "/test2.tf")
	var file *configs.File

	column := -1
	var diags hcl.Diagnostics
	if val, valDiags, isDot := helper.CheckAndGetConfig(parser, tempFile, vs.Position.Line+1, vs.Position.Character); isDot {
		diags = diags
		file = val
		column = vs.Position.Character - 1
	} else {
		diags = valDiags
		file = val
		column = vs.Position.Character
	}

	helper.DumpLog(diags)

	fileText, _ := ioutil.ReadFile(tempFile.Name())
	pos := helper.FindOffset(string(fileText), vs.Position.Line+1, column)

	//spew.Dump(err)

	var config *hclsyntax.Body

	for _, v := range file.ManagedResources {
		config = v.Config.(interface{}).(*hclsyntax.Body)
	}

	var result []lsp.CompletionItem

	if diags != nil || config == nil {
		helper.DumpLog(diags[0].Subject)
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        []lsp.CompletionItem{},
		}, nil
	}

	posHCL := hcl.Pos{
		Byte: pos,
	}

	expr := config.OutermostExprAtPos(posHCL)

	attr := config.AttributeAtPos(posHCL)
	helper.DumpLog(attr)
	helper.DumpLog(expr)

	var variables []hcl.Traversal

	if expr == nil {
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        []lsp.CompletionItem{},
		}, nil
	}

	origType := reflect.TypeOf(expr)

	if origType == hclstructs.BinaryOpExpr() {
		expr := expr.(*hclsyntax.BinaryOpExpr)
		if expr.LHS.Range().ContainsPos(posHCL) {
			variables = expr.LHS.Variables()
		} else if expr.RHS.Range().ContainsPos(posHCL) {
			variables = expr.RHS.Variables()
		}
	} else if origType == hclstructs.FunctionCallExpr() {
		expr := expr.(*hclsyntax.FunctionCallExpr)
		for _, arg := range expr.ExprCall().Arguments {
			if arg.Range().ContainsPos(posHCL) {
				variables = arg.Variables()
			}
		}
	} else {
		variables = expr.Variables()
	}
	//spew.Dump(config)
	if variables[0].RootName() == "var" {
		vars := variables[0]
		result = helper.ParseVariables(vars[1:], append(file.Variables, file2.Variables...), result)
	}

	//dumpLog(file2)
	//pos := hcl.Pos{
	//	Line:   vs.Position.Line + 1,
	//	Column: vs.Position.Character,
	//	Byte:   100,
	//}

	// //resourceType := file.ManagedResources[-1].Type
	// //file.ManagedResources[0].Config.Content(

	// provider, err := tfschema.NewClient(strings.Split(resourceType, "_")[0])
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //	//
	// provider_resource, err := provider.GetRawResourceTypeSchema(file.ManagedResources[0].Type)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //fmt.Println(provider_resource.Block)
	// //spew.Dump(file.ProviderConfigs)
	// //
	// res2 := provider_resource.Block.DecoderSpec()
	// res, diags := hcldec.Decode(file.ManagedResources[0].Config, res2, nil)

	// file.ManagedResources[0].Config.BlocksAtPos
	// dumpLog(res)
	// dumpLog(diags)
	return lsp.CompletionList{
		IsIncomplete: false,
		Items:        result,
	}, nil

}

func TextDocumentDidChange(ctx context.Context, vs lsp.DidChangeTextDocumentParams) error {
	tempFile.Truncate(0)
	tempFile.Seek(0, 0)
	tempFile.Write([]byte(vs.ContentChanges[0].Text))
	return nil
}

func TextDocumentDidOpen(ctx context.Context, vs lsp.DidOpenTextDocumentParams) error {
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

func main() {
	s := jrpc2.NewServer(handler.Map{
		"initialize":              handler.New(Initialize),
		"textDocument/completion": handler.New(TextDocumentComplete),
		"textDocument/didChange":  handler.New(TextDocumentDidChange),
		"textDocument/didOpen":    handler.New(TextDocumentDidOpen),
		"textDocument/didClose":   handler.New(TextDocumentDidClose),
		"exit":                    handler.New(Exit),
		"$/cancelRequest":         handler.New(CancelRequest),
	}, nil)

	f, err := os.OpenFile("tf-lsp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// Start the server on a channel comprising stdin/stdout.
	s.Start(channel.Header("")(os.Stdin, os.Stdout))
	log.Print("Server started")

	// Wait for the server to exit, and report any errors.
	if err := s.Wait(); err != nil {
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
