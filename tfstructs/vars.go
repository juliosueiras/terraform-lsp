package tfstructs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/lang"
	"github.com/juliosueiras/terraform-lsp/hclstructs"
	"github.com/juliosueiras/terraform-lsp/helper"

	//	"github.com/juliosueiras/terraform-lsp/loghelper"
	"os"
	"path/filepath"
	"reflect"

	"github.com/sourcegraph/go-lsp"
	"github.com/zclconf/go-cty/cty"
)

type GetVarAttributeRequest struct {
	Variables hcl.Traversal
	Result    []lsp.CompletionItem
	Files     *configs.Module
	Config    hcl.Body
	FileDir   string
}

func GetVarAttributeCompletion(request GetVarAttributeRequest) []lsp.CompletionItem {
	scope := lang.Scope{}

	targetDir := request.FileDir

	resultedDir := ""
	searchLevel := 4
	for dir := targetDir; dir != "" && searchLevel != 0; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, ".terraform")); err == nil {
			resultedDir = dir
			break
		}
		searchLevel -= 1
	}

	variables := map[string]cty.Value{
		"path": cty.ObjectVal(map[string]cty.Value{
			"cwd":    cty.StringVal(request.FileDir),
			"module": cty.StringVal(request.FileDir),
			"root":   cty.StringVal(resultedDir),
		}),
		"var":    cty.DynamicVal, // Need to check for undefined vars
		"module": cty.DynamicVal,
		"local":  cty.DynamicVal,
		"each":   cty.DynamicVal,
		"count":  cty.DynamicVal,
		"terraform": cty.ObjectVal(map[string]cty.Value{
			"workspace": cty.StringVal(""),
		}),
	}

	//rootName := ""
	if request.Variables.IsRelative() {
		request.Variables = request.Variables[1:]
	}

	if request.Variables.RootName() == "var" {
		vars := request.Variables

		request.Result = helper.ParseVariables(vars[1:], request.Files.Variables, request.Result)
	} else if request.Variables.RootName() == "local" {
		if len(request.Variables) > 1 {
			var found *configs.Local
			for _, v := range request.Files.Locals {
				if v.Name == request.Variables[1].(hcl.TraverseAttr).Name {
					found = v
					break
				}
			}

			if found != nil {
				testVal, _ := found.Expr.Value(
					&hcl.EvalContext{
						// Build Full Tree
						Variables: variables,
						Functions: scope.Functions(),
					},
				)

				origType := reflect.TypeOf(found.Expr)

				if origType == hclstructs.ObjectConsExpr() {
					items := found.Expr.(*hclsyntax.ObjectConsExpr).Items
					for _, v := range request.Variables[2:] {
						for _, l := range items {
							if v.(hcl.TraverseAttr).Name == l.KeyExpr.(*hclsyntax.ObjectConsKeyExpr).Wrapped.(*hclsyntax.ScopeTraversalExpr).AsTraversal().RootName() {
								origType2 := reflect.TypeOf(l.ValueExpr)

								if origType2 == hclstructs.ObjectConsExpr() {
									items = l.ValueExpr.(*hclsyntax.ObjectConsExpr).Items
								}
							}
						}
					}

					for _, v := range items {
						origType2 := reflect.TypeOf(v.ValueExpr)
						request.Result = append(request.Result, lsp.CompletionItem{
							Label:  v.KeyExpr.(*hclsyntax.ObjectConsKeyExpr).Wrapped.(*hclsyntax.ScopeTraversalExpr).AsTraversal().RootName(),
							Detail: fmt.Sprintf(" %s", hclstructs.GetExprStringType(origType2)),
						})
					}
				}

				request.Result = append(request.Result, helper.ParseOtherAttr(request.Variables[2:], testVal.Type(), request.Result)...)

				return request.Result
			}

		} else if len(request.Variables) == 1 {
			for _, v := range request.Files.Locals {
				origType := reflect.TypeOf(v.Expr)
				request.Result = append(request.Result, lsp.CompletionItem{
					Label:  v.Name,
					Detail: fmt.Sprintf(" local value(%s)", hclstructs.GetExprStringType(origType)),
				})
			}

			return request.Result
		} else {
			for _, v := range request.Files.Locals {
				origType := reflect.TypeOf(v.Expr)
				request.Result = append(request.Result, lsp.CompletionItem{
					Label:  v.Name,
					Detail: fmt.Sprintf(" local value(%s)", hclstructs.GetExprStringType(origType)),
				})
			}

			return request.Result
		}
	} else if request.Variables.RootName() == "data" {
		// Need refactoring
		if len(request.Variables) > 2 {
			re, _ := addrs.ParseAbsResourceInstanceStr(request.Variables.RootName() + "." + request.Variables[1].(hcl.TraverseAttr).Name)
			result := request.Files.ResourceByAddr(re.Resource.Resource)
			var providerType string
			if result != nil && result.ProviderConfigRef != nil {
				providerType = result.ProviderConfigRef.Name
			}
			res := GetDataSourceSchema(request.Variables[1].(hcl.TraverseAttr).Name, request.Config, request.FileDir, providerType)

			if res == nil {
				request.Result = append(request.Result, lsp.CompletionItem{
					Label:  "",
					Detail: " No such data source",
				})
				return request.Result
			}

			request.Result = helper.ParseOtherAttr(request.Variables[3:], res.Schema.Block.ImpliedType(), request.Result)
			return request.Result
		} else if len(request.Variables) == 2 {
			for _, v := range request.Files.DataResources {
				if v.Type == request.Variables[1].(hcl.TraverseAttr).Name {
					request.Result = append(request.Result, lsp.CompletionItem{
						Label:  v.Name,
						Detail: " data source instance",
					})
				}
			}
			return request.Result
		} else {
			for _, v := range request.Files.DataResources {
				existed := false
				for _, e := range request.Result {
					if e.Label == v.Type {
						existed = true
						break
					}
				}

				if !existed {
					request.Result = append(request.Result, lsp.CompletionItem{
						Label:  v.Type,
						Detail: " data resource",
					})
				}
			}

			return request.Result
		}
	} else {
		if len(request.Variables) > 1 {
			re, _ := addrs.ParseAbsResourceInstanceStr(request.Variables.RootName() + "." + request.Variables[1].(hcl.TraverseAttr).Name)
			result := request.Files.ResourceByAddr(re.Resource.Resource)
			var providerType string
			if result != nil && result.ProviderConfigRef != nil {
				providerType = result.ProviderConfigRef.Name
			}

			res := GetResourceSchema(request.Variables.RootName(), request.Config, request.FileDir, providerType)

			if res == nil {
				request.Result = append(request.Result, lsp.CompletionItem{
					Label:  "",
					Detail: " No such resource",
				})
				return request.Result
			}

			request.Result = helper.ParseOtherAttr(request.Variables[2:], res.Schema.Block.ImpliedType(), request.Result)
			return request.Result
		} else {
			for _, v := range request.Files.ManagedResources {
				if v.Type == request.Variables.RootName() {
					request.Result = append(request.Result, lsp.CompletionItem{
						Label:  v.Name,
						Detail: " resource instance",
					})
				}
			}
		}
	}

	return request.Result
}
