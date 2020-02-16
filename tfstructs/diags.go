package tfstructs

import (
	"fmt"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
	//"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/juliosueiras/terraform-lsp/memfs"
	"github.com/sourcegraph/go-lsp"
	"github.com/spf13/afero"
	"path/filepath"
)

func GetDiagnostics(fileName string, originalFile string) []lsp.Diagnostic {
  
	parser := configs.NewParser(memfs.MemFs)
	result := make([]lsp.Diagnostic, 0)
	originalFileName := originalFile
  
	if exist, _ := afero.Exists(memfs.MemFs, fileName); !exist {
		return result
	}

	if exist, _ := afero.Exists(memfs.MemFs, originalFile); !exist {
		originalFile = fileName
	}

	_, hclDiags := parser.LoadHCLFile(fileName)

	for _, diag := range hclDiags {
		result = append(result, lsp.Diagnostic{
			Severity: lsp.DiagnosticSeverity(diag.Severity),
			Message:  diag.Detail,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      diag.Subject.Start.Line - 1,
					Character: diag.Subject.Start.Column - 1,
				},
				End: lsp.Position{
					Line:      diag.Subject.End.Line - 1,
					Character: diag.Subject.End.Column - 1,
				},
			},
			Source: "HCL",
		})
	}

	cfg, tfDiags := parser.LoadConfigFile(fileName)
	parser.ForceFileSource(originalFileName, []byte(""))
	extra, _ := parser.LoadConfigDir(filepath.Dir(originalFileName))

	resourceTypes := map[string]map[string]cty.Value{}

	if extra != nil {
		for _, v := range extra.ManagedResources {
			if resourceTypes[v.Type] == nil {
				resourceTypes[v.Type] = map[string]cty.Value{}
			}

			resourceTypes[v.Type][v.Name] = cty.DynamicVal
		}
	}

	for _, v := range cfg.ManagedResources {
		if resourceTypes[v.Type] == nil {
			resourceTypes[v.Type] = map[string]cty.Value{}
		}

		resourceTypes[v.Type][v.Name] = cty.DynamicVal
	}

	variables := map[string]cty.Value{
		"path": cty.ObjectVal(map[string]cty.Value{
			"cwd":    cty.StringVal(""),
			"module": cty.StringVal(""),
		}),
		"var":    cty.DynamicVal, // Need to check for undefined vars
		"module": cty.DynamicVal,
		"local":  cty.DynamicVal,
		"each":   cty.DynamicVal,
		"count":  cty.DynamicVal,
	}

	for k, v := range resourceTypes {
		variables[k] = cty.ObjectVal(v)
	}

	dataTypes := map[string]map[string]cty.Value{}

	if extra != nil {
		for _, v := range extra.DataResources {
			if dataTypes[v.Type] == nil {
				dataTypes[v.Type] = map[string]cty.Value{}
			}

			dataTypes[v.Type][v.Name] = cty.DynamicVal
		}
	}

	for _, v := range cfg.DataResources {
		if dataTypes[v.Type] == nil {
			dataTypes[v.Type] = map[string]cty.Value{}
		}

		dataTypes[v.Type][v.Name] = cty.DynamicVal
	}

	resultDataTypes := map[string]cty.Value{}

	for k, v := range dataTypes {
		resultDataTypes[k] = cty.ObjectVal(v)
	}

	variables["data"] = cty.ObjectVal(resultDataTypes)

	for _, diag := range tfDiags {
		result = append(result, lsp.Diagnostic{
			Severity: lsp.DiagnosticSeverity(diag.Severity),
			Message:  diag.Detail,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      diag.Subject.Start.Line - 1,
					Character: diag.Subject.Start.Column - 1,
				},
				End: lsp.Position{
					Line:      diag.Subject.End.Line - 1,
					Character: diag.Subject.End.Column - 1,
				},
			},
			Source: "Terraform",
		})
	}
	//	cfg, diags := configload.NewLoader(&configload.Config{
	//		ModulesDir: ".terraform/modules",
	//	})
	//	helper.DumpLog(diags)
	//	config, diags2 := cfg.LoadConfig(filepath.Dir(originalFile))
	//	input := addrs.InputVariable{
	//		Name: "test_attr2",
	//	}
	// helper.DumpLog(cfg2.ModuleCalls[0].Config.(*hclsyntax.Body).Attributes["test_attr2"].Expr)
	//	evalModule := terraform.EvalModuleCallArgument{
	//		Addr:   input,
	//		Config: config.Children["test"].Module.Variables["test_attr2"],
	//		Values: make(map[string]cty.Value),
	//		Expr:   cfg2.ModuleCalls[0].Config.(*hclsyntax.Body).Attributes["test_attr2"].Expr,
	//	}
	//
	//	result2, diags3 := evalModule.Eval(&terraform.MockEvalContext{
	//		PathPath: addrs.RootModuleInstance,
	//	})
	//
	//	helper.DumpLog(diags2)
	//	helper.DumpLog(diags3)
	//	helper.DumpLog(result2)
	//	helper.DumpLog(config)
	for _, v := range cfg.ProviderConfigs {
		providerType := v.Name

		tfSchema := GetProviderSchemaForDiags(providerType, v.Config, filepath.Dir(originalFile), variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      diag.Subject.Start.Line - 1,
							Character: diag.Subject.Start.Column - 1,
						},
						End: lsp.Position{
							Line:      diag.Subject.End.Line - 1,
							Character: diag.Subject.End.Column - 1,
						},
					},
					Source: "Terraform Schema",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Provider %s does not exist", v.Name),
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      v.NameRange.Start.Line - 1,
						Character: v.NameRange.Start.Column - 1,
					},
					End: lsp.Position{
						Line:      v.NameRange.End.Line - 1,
						Character: v.NameRange.End.Column - 1,
					},
				},
				Source: "Terraform Schema",
			})
		}
	}

	for _, v := range cfg.ManagedResources {
		resourceType := v.Type

		var providerType string
		if v.ProviderConfigRef != nil {
			providerType = v.ProviderConfigRef.Name
		}

		tfSchema := GetResourceSchemaForDiags(resourceType, v.Config, filepath.Dir(originalFile), providerType, variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      diag.Subject.Start.Line - 1,
							Character: diag.Subject.Start.Column - 1,
						},
						End: lsp.Position{
							Line:      diag.Subject.End.Line - 1,
							Character: diag.Subject.End.Column - 1,
						},
					},
					Source: "Terraform Schema",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Resource %s does not exist", v.Type),
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      v.TypeRange.Start.Line - 1,
						Character: v.TypeRange.Start.Column - 1,
					},
					End: lsp.Position{
						Line:      v.TypeRange.End.Line - 1,
						Character: v.TypeRange.End.Column - 1,
					},
				},
				Source: "Terraform Schema",
			})
		}
	}

	for _, v := range cfg.DataResources {
		resourceType := v.Type
		var providerType string
		if v.ProviderConfigRef != nil {
			providerType = v.ProviderConfigRef.Name
		}

		tfSchema := GetDataSourceSchemaForDiags(resourceType, v.Config, filepath.Dir(originalFile), providerType, variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      diag.Subject.Start.Line - 1,
							Character: diag.Subject.Start.Column - 1,
						},
						End: lsp.Position{
							Line:      diag.Subject.End.Line - 1,
							Character: diag.Subject.End.Column - 1,
						},
					},
					Source: "Terraform",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Data source %s does not exist", v.Type),
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      v.TypeRange.Start.Line - 1,
						Character: v.TypeRange.Start.Column - 1,
					},
					End: lsp.Position{
						Line:      v.TypeRange.End.Line - 1,
						Character: v.TypeRange.End.Column - 1,
					},
				},
				Source: "Terraform",
			})
		}
	}
	//
	//	//spew.Dump(file.ManagedResources[0].Config.Content(nil))

	return result
}
