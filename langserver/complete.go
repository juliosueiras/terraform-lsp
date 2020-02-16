package langserver

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/lang"
	lsp "github.com/sourcegraph/go-lsp"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	//"github.com/hashicorp/terraform/terraform"
	"github.com/juliosueiras/terraform-lsp/hclstructs"
	"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/juliosueiras/terraform-lsp/tfstructs"
	"github.com/juliosueiras/terraform-lsp/memfs"
  "github.com/spf13/afero"
)

func TextDocumentComplete(ctx context.Context, vs lsp.CompletionParams) (lsp.CompletionList, error) {
  //log.Println(tfstructs.Clients)
	parser := configs.NewParser(memfs.MemFs)

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
	var haveDot bool
	file, diags, column, hclFile, haveDot = helper.CheckAndGetConfig(parser, tempFile, vs.Position.Line+1, vs.Position.Character)

	resultFiles = append(resultFiles, file)

	files, diags := configs.NewModule(resultFiles, nil)

	fileText, _ := afero.ReadFile(memfs.MemFs, tempFile.Name())
	pos := helper.FindOffset(string(fileText), vs.Position.Line+1, column)

	var result []lsp.CompletionItem

	posHCL := hcl.Pos{
		Byte: pos,
	}

	var extraProvider string
	if files.ProviderConfigs != nil {
		for k, _ := range files.ProviderConfigs {
			if k == "google-beta" {
				extraProvider = "google-beta"
			}
		}
	}

	if r, found, _ := tfstructs.GetTypeCompletion(result, fileDir, hclFile, posHCL, extraProvider); found {
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

	//if expr != nil {
	//	scope := lang.Scope{}
	//	s, w := scope.EvalExpr(expr, cty.DynamicPseudoType)
	//	helper.DumpLog(s)
	//	helper.DumpLog(w)
	//}

	if expr == nil && attr == nil && blocks == nil {
		attrs, _ := config.JustAttributes()
		fileText, _ := afero.ReadFile(memfs.MemFs, tempFile.Name())
		pos := helper.FindOffset(string(fileText), vs.Position.Line+1, column+1)

		posHCL2 := hcl.Pos{
			Byte: pos,
		}

		for _, v := range attrs {
			origType := reflect.TypeOf(v.Expr)
			if origType == hclstructs.LiteralValueExpr() {
				if v.Expr.(*hclsyntax.LiteralValueExpr).Range().ContainsPos(posHCL2) {
					scope := lang.Scope{}

					// Add Detaults
					defaults := map[string]string{
						"local":     " locals",
						"path":      " path",
						"terraform": " workspace",
						"var":       " variable",
						"module":    " module",
						"data":      " data source",
					}

					for k, v := range defaults {
						result = append(result, lsp.CompletionItem{
							Label:  k,
							Detail: v,
						})
					}
					for k, v := range scope.Functions() {
						var params []string

						for _, x := range v.Params() {
							params = append(params, x.Name)
						}

						result = append(result, lsp.CompletionItem{
							Label:      fmt.Sprintf("%s(%s)", k, strings.Join(params, ",")),
							InsertText: k,
							Detail:     " function",
						})

					}

					for _, v := range files.ManagedResources {
						existed := false
						for _, e := range result {
							if e.Label == v.Type {
								existed = true
								break
							}
						}

						if !existed {
							result = append(result, lsp.CompletionItem{
								Label:  v.Type,
								Detail: " resource",
							})
						}
					}

					return lsp.CompletionList{
						IsIncomplete: false,
						Items:        result,
					}, nil

				}
			}
		}

		//hclsyntax.LiteralValueExpr
		if r, found, _ := tfstructs.GetAttributeCompletion(result, configType, origConfig, fileDir); found {
			return r, nil
		}
	}

	// Block is Block, not resources
	//test := config.BlocksAtPos(posHCL)
	if blocks != nil && attr == nil {
		//helper.DumpLog(blocks)
		if blocks[0].Type == "provisioner" {
			helper.DumpLog(blocks)
			if len(blocks) == 1 {

				if r, found, _ := tfstructs.GetAttributeCompletion(result, "provisioner", blocks[0], fileDir); found {
					return r, nil
				}
			} else {
				if r, found, _ := tfstructs.GetNestingCompletion(blocks[1:], result, "provisioner", blocks[0], fileDir); found {
					return r, nil
				}
			}

		} else if blocks[0].Type == "dynamic" {
			if len(blocks) == 1 {
				result = append(result, lsp.CompletionItem{
					Label:  "for_each",
					Detail: " dynamic",
				})
				return lsp.CompletionList{
					IsIncomplete: false,
					Items:        result,
				}, nil
			}

			dynamicBlock := blocks[0]
			blocks := blocks[1:]
			blocks[0].Type = dynamicBlock.Labels[0]
			if r, found, _ := tfstructs.GetNestingCompletion(blocks, result, configType, origConfig, fileDir); found {
				return r, nil
			}
		}
		if r, found, _ := tfstructs.GetNestingCompletion(blocks, result, configType, origConfig, fileDir); found {
			return r, nil
		}
	}

	if expr != nil {
		helper.DumpLog("Found Expression")
		helper.DumpLog(expr)
		//.*for.*in\s+([^:]*)
		//te, te2 := hclsyntax.ParseExpression([]byte("aws[0].test"), "test", hcl.Pos{
		//	Line:   0,
		//	Column: 0,
		//})
		//helper.DumpLog(te)
		//helper.DumpLog(te2)
		origType := reflect.TypeOf(expr)
		if origType == hclstructs.LiteralValueExpr() {
			if expr.(*hclsyntax.LiteralValueExpr).Val.Type().HasDynamicTypes() {

				textLines := strings.Split(string(fileText), "\n")

				re := regexp.MustCompile(".*for.*in\\s+([^:]*)")
				searchResult := re.FindSubmatch([]byte(textLines[vs.Position.Line]))

				if searchResult != nil {
					helper.DumpLog(searchResult[1])
					dynamicExpr, _ := hclsyntax.ParseExpression([]byte(searchResult[1]), "test", hcl.Pos{
						Line:   0,
						Column: 0,
					})

					if len(dynamicExpr.Variables()) != 0 {
						result = tfstructs.GetVarAttributeCompletion(tfstructs.GetVarAttributeRequest{
							Variables: dynamicExpr.Variables()[0],
							Result:    result,
							Files:     files,
							Config:    config,
							FileDir:   fileDir,
						})
						return lsp.CompletionList{
							IsIncomplete: false,
							Items:        result,
						}, nil
					}
				}
			}
		}
		//reflect.New(origType)
		if origType == hclstructs.ForExpr() {
			expr := expr.(*hclsyntax.ForExpr)
			helper.DumpLog(expr)
			resultName := []string{}
			helper.DumpLog(expr.ValExpr.Range().ContainsPos(posHCL))
			if expr.ValExpr.Range().ContainsPos(posHCL) {
				if reflect.TypeOf(expr.CollExpr) == hclstructs.ScopeTraversalExpr() {
					resultName = append(resultName, expr.CollExpr.(*hclsyntax.ScopeTraversalExpr).AsTraversal().RootName())

					for _, v := range expr.CollExpr.(*hclsyntax.ScopeTraversalExpr).AsTraversal()[1:] {
						if reflect.TypeOf(v) == hclstructs.TraverseAttr() {
							resultName = append(resultName, v.(hcl.TraverseAttr).Name)
						} else if reflect.TypeOf(v) == hclstructs.TraverseIndex() {
							resultName = append(resultName, "<index>")
						}
					}
				}

				scopeExpr := expr.ValExpr.(*hclsyntax.ScopeTraversalExpr)
				helper.DumpLog(haveDot)
				helper.DumpLog((len(scopeExpr.AsTraversal()) == 1 && haveDot) || len(scopeExpr.AsTraversal()) > 1)
				if len(scopeExpr.AsTraversal()) == 1 && !haveDot {
					helper.DumpLog(vs.Position.Character)
					result = append(result, lsp.CompletionItem{
						Label:  expr.ValVar,
						Detail: fmt.Sprintf(" foreach var(%s)", strings.Join(resultName, ".")),
					})
				} else if (len(scopeExpr.AsTraversal()) == 1 && haveDot) || len(scopeExpr.AsTraversal()) > 1 {
					forVars := expr.CollExpr.(*hclsyntax.ScopeTraversalExpr).AsTraversal()
					for _, v := range scopeExpr.Traversal[1:] {
						forVars = append(forVars, v)
					}
					result = tfstructs.GetVarAttributeCompletion(tfstructs.GetVarAttributeRequest{
						Variables: forVars,
						Result:    result,
						Files:     files,
						Config:    config,
						FileDir:   fileDir,
					})
					return lsp.CompletionList{
						IsIncomplete: false,
						Items:        result,
					}, nil

				}

				return lsp.CompletionList{
					IsIncomplete: false,
					Items:        result,
				}, nil
			}
		}
		//tests, errxs := lang.ReferencesInExpr(expr)
		if origType != hclstructs.ObjectConsExpr() {
			variables := hclstructs.GetExprVariables(origType, expr, posHCL)

			if len(variables) != 0 {
				result = tfstructs.GetVarAttributeCompletion(tfstructs.GetVarAttributeRequest{
					Variables: variables[0],
					Result:    result,
					Files:     files,
					Config:    config,
					FileDir:   fileDir,
				})
				return lsp.CompletionList{
					IsIncomplete: false,
					Items:        result,
				}, nil

			} else {
				scope := lang.Scope{}

				// Add Detaults
				defaults := map[string]string{
					"local":     " locals",
					"path":      " path",
					"terraform": " workspace",
					"var":       " variable",
					"module":    " module",
					"data":      " data source",
				}
				for k, v := range defaults {
					result = append(result, lsp.CompletionItem{
						Label:  k,
						Detail: v,
					})
				}
				for k, v := range scope.Functions() {
					var params []string

					for _, x := range v.Params() {
						params = append(params, x.Name)
					}

					result = append(result, lsp.CompletionItem{
						Label:      fmt.Sprintf("%s(%s)", k, strings.Join(params, ",")),
						InsertText: k,
						Detail:     " function",
					})

				}

				for _, v := range files.ManagedResources {
					existed := false
					for _, e := range result {
						if e.Label == v.Type {
							existed = true
							break
						}
					}

					if !existed {
						result = append(result, lsp.CompletionItem{
							Label:  v.Type,
							Detail: " resource",
						})
					}
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
