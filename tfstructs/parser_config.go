package tfstructs

import (
	"fmt"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/juliosueiras/terraform-lsp/hclstructs"
	"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/sourcegraph/go-lsp"
	"github.com/zclconf/go-cty/cty"
	"strings"
)

func GetNestingAttributeCompletion(attr *hcl.Attribute, result []lsp.CompletionItem, configType string, origConfig interface{}, fileDir string, posHCL hcl.Pos) (lsp.CompletionList, bool, error) {

	topName := attr.Name

	res := hclstructs.GetExprVariables(hclstructs.ObjectConsExpr(), attr.Expr, posHCL)

	attrs := hcl.Traversal{}
	if len(res) != 0 {
		attrs = res[0]
	}
	resultTraversal := hcl.Traversal{
		hcl.TraverseAttr{
			Name: topName,
		},
	}

	resultTraversal = append(resultTraversal, attrs...)
	switch configType {
	case "module":
		if moduleVars, found := GetModuleVariables(origConfig.(*configs.ModuleCall).SourceAddr, origConfig.(*configs.ModuleCall).Config, fileDir); found {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        helper.ParseVariables(resultTraversal, moduleVars, result),
			}, true, nil
		}
	case "resource":
	case "data":
		return lsp.CompletionList{
			IsIncomplete: false,
			//Items:        helper.ParseVariables(resultTraversal, origConfig.(*configs.Resource)., result),
			Items: result,
		}, true, nil

	}
	return lsp.CompletionList{
		IsIncomplete: false,
		Items:        result,
	}, true, nil

	//case "module":
	//	origConfig := origConfig.(*configs.ModuleCall)

	//	if res, found := GetModuleVariables(origConfig.SourceAddr, origConfig.Config, fileDir); found {
	//		for k, v := range res {
	//			result = append(result, lsp.CompletionItem{
	//				Label:         k,
	//				Detail:        fmt.Sprintf("  %s", v.Type.FriendlyName()),
	//				Documentation: v.Description,
	//			})
	//		}
	//		return lsp.CompletionList{
	//			IsIncomplete: false,
	//			Items:        result,
	//		}, true, nil

	//	}
	//	return lsp.CompletionList{
	//		IsIncomplete: false,
	//		Items:        result,
	//	}, true, nil
}

func GetTypeCompletion(result []lsp.CompletionItem, fileDir string, hclFile *hclsyntax.Body, posHCL hcl.Pos, extraProvider string) (lsp.CompletionList, bool, error) {
	for _, v := range hclFile.Blocks {
		blockType := v.Type
		for i, r := range v.LabelRanges {
			if r.ContainsPos(posHCL) {
				if len(strings.Split(v.Labels[i], "_")) == 1 {
					for providerName, providerInfo := range OfficialProviders {
						result = append(result, lsp.CompletionItem{
							Label:  providerName,
							Detail: fmt.Sprintf(" %s (%s)", providerInfo.Type, providerInfo.Name),
						})
					}

					return lsp.CompletionList{
						IsIncomplete: true,
						Items:        result,
					}, true, nil
				}

				if blockType != "resource" && blockType != "data" {
					return lsp.CompletionList{
						IsIncomplete: true,
						Items:        result,
					}, true, nil
				}

				var test *Client
				var includeExtraProvider string
				if strings.Contains(v.Labels[i], "google") && extraProvider == "google-beta" {
					test, _ = GetProvider(extraProvider, fileDir)
					includeExtraProvider = "(include Beta)"
				} else {
					test, _ = GetProvider(v.Labels[i], fileDir)
				}
				if test == nil {
					result = append(result, lsp.CompletionItem{
						Label:  v.Labels[i],
						Detail: " Did You Forgot to do terraform init?",
					})
					return lsp.CompletionList{
						IsIncomplete: false,
						Items:        result,
					}, true, nil
				}
				defer test.Kill()

				var res []string
				var resultType string

				switch v.Type {
				case "resource":
					resultType = "resource"
					res, _ = test.GetResourceTypes()
				case "data":
					resultType = "data source"
					res, _ = test.GetDataSourceTypes()
				}

				if res != nil {
					for _, resource := range res {
						if strings.HasPrefix(resource, v.Labels[i]) {
							result = append(result, lsp.CompletionItem{
								Label:  resource,
								Detail: fmt.Sprintf(" %s %s", resultType, includeExtraProvider),
							})
						}
					}

					return lsp.CompletionList{
						IsIncomplete: true,
						Items:        result,
					}, true, nil
				}
			}
		}
	}

	return lsp.CompletionList{}, false, nil
}

func GetConfig(file *configs.File, posHCL hcl.Pos) (*hclsyntax.Body, interface{}, string) {
	var config *hclsyntax.Body
	var origConfig interface{}
	var configType string

	//*configs.Resource
	for _, v := range file.ManagedResources {
		if v.Config.(*hclsyntax.Body).Range().ContainsPos(posHCL) {
			configType = "resource"
			origConfig = v
			config = v.Config.(interface{}).(*hclsyntax.Body)
		}
	}

	//* configs.Provider
	for _, v := range file.ProviderConfigs {
		if v.Config.(*hclsyntax.Body).Range().ContainsPos(posHCL) {
			configType = "provider"
			origConfig = v
			config = v.Config.(interface{}).(*hclsyntax.Body)
		}
	}

	//* configs.Resource
	for _, v := range file.DataResources {
		if v.Config.(*hclsyntax.Body).Range().ContainsPos(posHCL) {
			configType = "data"
			origConfig = v
			config = v.Config.(interface{}).(*hclsyntax.Body)
		}
	}

	//* configs.Backends
	for _, v := range file.Backends {
		helper.DumpLog(v)
		helper.DumpLog(TerraformBackends["remote"])
		//		helper.DumpLog(backend.Backend{
		//
		//		})
		//		if v.Config.(*hclsyntax.Body).Range().ContainsPos(posHCL) {
		//			configType = "backend"
		//			origConfig = v
		//			config = v.Config.(interface{}).(*hclsyntax.Body)
		//		}
	}

	for _, v := range file.ModuleCalls {
		if v.Config.(*hclsyntax.Body).Range().ContainsPos(posHCL) {
			configType = "module"
			origConfig = v
			config = v.Config.(interface{}).(*hclsyntax.Body)
		}
	}

	return config, origConfig, configType
}

func GetAttributeCompletion(result []lsp.CompletionItem, configType string, origConfig interface{}, fileDir string) (lsp.CompletionList, bool, error) {
	switch configType {
	case "provider":
		origConfig := origConfig.(*configs.Provider)

		res, _ := GetProvider(origConfig.Name, fileDir)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		defer res.Kill()
		schema, _ := res.GetRawProviderSchema()

		for k, v := range schema.Block.Attributes {
			if v.Optional || v.Required {
				result = append(result, lsp.CompletionItem{
					Label:         k,
					Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
					Documentation: v.Description,
				})
			}
		}

		for p, v := range schema.Block.BlockTypes {
			result = append(result, lsp.CompletionItem{
				Label:  p,
				Detail: " " + v.Nesting.String(),
			})
		}
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        result,
		}, true, nil

	case "resource":
		origConfig := origConfig.(*configs.Resource)
		var providerType string
		if origConfig.ProviderConfigRef != nil {
			providerType = origConfig.ProviderConfigRef.Name
		}

		res := GetResourceSchema(origConfig.Type, origConfig.Config, fileDir, providerType)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		helper.DumpLog(res)

		for k, v := range res.Schema.Block.Attributes {
			if v.Optional || v.Required {
				result = append(result, lsp.CompletionItem{
					Label:         k,
					Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
					Documentation: v.Description,
				})
			}
		}

		for p, v := range res.Schema.Block.BlockTypes {
			result = append(result, lsp.CompletionItem{
				Label:  p,
				Detail: " " + v.Nesting.String(),
			})
		}
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        result,
		}, true, nil
	case "data":
		origConfig := origConfig.(*configs.Resource)

		var providerType string
		if origConfig.ProviderConfigRef != nil {
			providerType = origConfig.ProviderConfigRef.Name
		}

		res := GetDataSourceSchema(origConfig.Type, origConfig.Config, fileDir, providerType)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		for k, v := range res.Schema.Block.Attributes {
			if v.Optional || v.Required {
				result = append(result, lsp.CompletionItem{
					Label:         k,
					Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
					Documentation: v.Description,
				})
			}
		}

		for p, v := range res.Schema.Block.BlockTypes {
			result = append(result, lsp.CompletionItem{
				Label:  p,
				Detail: " " + v.Nesting.String(),
			})
		}
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        result,
		}, true, nil
	case "module":
		origConfig := origConfig.(*configs.ModuleCall)

		if res, found := GetModuleVariables(origConfig.SourceAddr, origConfig.Config, fileDir); found {
			for k, v := range res {
				result = append(result, lsp.CompletionItem{
					Label:         k,
					Detail:        fmt.Sprintf("  %s", v.Type.FriendlyName()),
					Documentation: v.Description,
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil

		}
		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        result,
		}, true, nil

	}

	return lsp.CompletionList{}, false, nil
}

func checkRequire(v *configschema.Attribute) string {
	if v.Required {
		return "Required"
	} else {
		return "Optional"
	}
}

func GetNestingCompletion(blocks []*hcl.Block, result []lsp.CompletionItem, configType string, origConfig interface{}, fileDir string) (lsp.CompletionList, bool, error) {
	switch configType {
	case "provider":
		origConfig := origConfig.(*configs.Provider)
		res, _ := GetProvider(origConfig.Name, fileDir)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}
		defer res.Kill()
		schema, _ := res.GetRawProviderSchema()

		var resultBlock *configschema.NestedBlock
		searchBlockTypes := schema.Block.BlockTypes
		for _, block := range blocks {
			if searchBlockTypes[block.Type] != nil {
				resultBlock = searchBlockTypes[block.Type]
				searchBlockTypes = searchBlockTypes[block.Type].BlockTypes
			}
		}

		if resultBlock != nil {
			for k, v := range resultBlock.Attributes {
				if v.Optional || v.Required {
					result = append(result, lsp.CompletionItem{
						Label:         k,
						Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
						Documentation: v.Description,
					})
				}
			}

			for p, v := range resultBlock.BlockTypes {
				result = append(result, lsp.CompletionItem{
					Label:  p,
					Detail: " " + v.Nesting.String(),
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

	case "resource":
		origConfig := origConfig.(*configs.Resource)
		var providerType string
		if origConfig.ProviderConfigRef != nil {
			providerType = origConfig.ProviderConfigRef.Name
		}
		res := GetResourceSchema(origConfig.Type, origConfig.Config, fileDir, providerType)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		var resultBlock *configschema.NestedBlock
		var resultType interface{}
		searchBlockTypes := res.Schema.Block.BlockTypes
		searchAttributes := res.Schema.Block.Attributes

		for _, block := range blocks {
			if searchBlockTypes[block.Type] != nil {
				resultBlock = searchBlockTypes[block.Type]
				resultType = nil
				searchBlockTypes = searchBlockTypes[block.Type].BlockTypes
			}

			if searchAttributes[block.Type] != nil {
				if searchAttributes[block.Type].Type.IsSetType() {
					if searchAttributes[block.Type].Type.SetElementType().IsObjectType() {
						resultBlock = nil
						resultType = searchAttributes[block.Type].Type
					}
				}
			}
		}

		if resultType != nil {
			for k, v := range resultType.(cty.Type).SetElementType().AttributeTypes() {
				result = append(result, lsp.CompletionItem{
					Label:  k,
					Detail: " " + v.FriendlyName(),
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		if resultBlock != nil {

			for k, v := range resultBlock.Attributes {
				if v.Optional || v.Required {
					result = append(result, lsp.CompletionItem{
						Label:         k,
						Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
						Documentation: v.Description,
					})
				}
			}

			for p, v := range resultBlock.BlockTypes {
				result = append(result, lsp.CompletionItem{
					Label:  p,
					Detail: " " + v.Nesting.String(),
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}
	case "data":
		origConfig := origConfig.(*configs.Resource)
		var providerType string
		if origConfig.ProviderConfigRef != nil {
			providerType = origConfig.ProviderConfigRef.Name
		}

		res := GetDataSourceSchema(origConfig.Type, origConfig.Config, fileDir, providerType)
		if res == nil {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		var resultBlock *configschema.NestedBlock
		var resultType interface{}
		searchBlockTypes := res.Schema.Block.BlockTypes
		searchAttributes := res.Schema.Block.Attributes

		for _, block := range blocks {
			if searchBlockTypes[block.Type] != nil {
				resultBlock = searchBlockTypes[block.Type]
				resultType = nil
				searchBlockTypes = searchBlockTypes[block.Type].BlockTypes
			}

			if searchAttributes[block.Type] != nil {
				if searchAttributes[block.Type].Type.IsSetType() {
					if searchAttributes[block.Type].Type.SetElementType().IsObjectType() {
						resultBlock = nil
						resultType = searchAttributes[block.Type].Type
					}
				}
			}
		}

		if resultType != nil {
			for k, v := range resultType.(cty.Type).SetElementType().AttributeTypes() {
				result = append(result, lsp.CompletionItem{
					Label:  k,
					Detail: " " + v.FriendlyName(),
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		if resultBlock != nil {

			for k, v := range resultBlock.Attributes {
				if v.Optional || v.Required {
					result = append(result, lsp.CompletionItem{
						Label:         k,
						Detail:        fmt.Sprintf(" (%s) %s", checkRequire(v), v.Type.FriendlyName()),
						Documentation: v.Description,
					})
				}
			}

			for p, v := range resultBlock.BlockTypes {
				result = append(result, lsp.CompletionItem{
					Label:  p,
					Detail: " " + v.Nesting.String(),
				})
			}
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}
	case "module":
		origConfig := origConfig.(*configs.ModuleCall)
		res, found := GetModuleVariables(origConfig.SourceAddr, origConfig.Config, fileDir)
		if !found {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		var resultTraversal hcl.Traversal

		for _, block := range blocks {
			resultTraversal = append(resultTraversal, hcl.TraverseAttr{
				Name: block.Type,
			})
		}
		tempResult := []lsp.CompletionItem{}

		resultVars := helper.ParseVariables(resultTraversal, res, tempResult)

		if len(resultVars) == 0 {
			return lsp.CompletionList{
				IsIncomplete: false,
				Items:        result,
			}, true, nil
		}

		result = append(result, resultVars...)

		return lsp.CompletionList{
			IsIncomplete: false,
			Items:        result,
		}, true, nil
	}
	return lsp.CompletionList{}, false, nil
}

func GetTopLevelCompletion() []lsp.CompletionItem {
	return []lsp.CompletionItem{
		lsp.CompletionItem{
			Label:  "resource",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "data",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "module",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "output",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "variable",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "provider",
			Detail: " type",
		},
		lsp.CompletionItem{
			Label:  "terraform",
			Detail: " type",
		},
	}

}
