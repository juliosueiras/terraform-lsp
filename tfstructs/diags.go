package tfstructs

import (
	"fmt"
	"github.com/hashicorp/terraform/configs"
	"github.com/sourcegraph/go-lsp"
	"os"
	"path/filepath"
)

func GetDiagnostics(fileName string, originalFile string) []lsp.Diagnostic {
	parser := configs.NewParser(nil)
	result := make([]lsp.Diagnostic, 0)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return result
	}

	if _, err := os.Stat(originalFile); os.IsNotExist(err) {
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

		tfSchema := GetProviderSchema(providerType, v.Config, filepath.Dir(originalFile))

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

		tfSchema := GetResourceSchema(resourceType, v.Config, filepath.Dir(originalFile))

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

		tfSchema := GetDataSourceSchema(resourceType, v.Config, filepath.Dir(originalFile))

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
