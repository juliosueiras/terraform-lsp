package tfstructs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juliosueiras/terraform-lsp/memfs"
	// "github.com/juliosueiras/terraform-lsp/helper"

	v2 "github.com/hashicorp/hcl/v2"
	oldHCL2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/configs"

	terragruntConfig "github.com/gruntwork-io/terragrunt/config"
	terragruntOptions "github.com/gruntwork-io/terragrunt/options"

	"github.com/sourcegraph/go-lsp"
	"github.com/spf13/afero"
	"github.com/zclconf/go-cty/cty"
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

	var hclDiags v2.Diagnostics
	isTFVars := (filepath.Ext(originalFile) == ".tfvars")
	isTerragrunt := (filepath.Base(originalFile) == "terragrunt.hcl")

	var diagName string

	if isTFVars {
		_, hclDiags = parser.LoadValuesFile(fileName)
		diagName = "TFVars"
	} else if isTerragrunt {
		fileContent, _ := afero.ReadFile(memfs.MemFs, fileName)

		_, terragruntDiags := terragruntConfig.ParseConfigString(string(fileContent), &terragruntOptions.TerragruntOptions{}, &terragruntConfig.IncludeConfig{}, originalFile)

		if terragruntDiags == nil {
			return result
		}

		for _, diag := range terragruntDiags.(oldHCL2.Diagnostics) {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(diag.Severity),
				Message:  diag.Detail,
				Range:    rangeOf(castToV2Range(*diag.Subject)),
				Source:   "Terragrunt",
			})
		}

		return result

	} else {
		_, hclDiags = parser.LoadHCLFile(fileName)
		diagName = "HCL"
	}

	for _, diag := range hclDiags {
		result = append(result, lsp.Diagnostic{
			Severity: lsp.DiagnosticSeverity(diag.Severity),
			Message:  diag.Detail,
			Range:    rangeOf(*diag.Subject),
			Source:   diagName,
		})
	}

	if isTFVars {
		return result
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

	targetDir := filepath.Dir(originalFileName)

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
			"cwd":    cty.StringVal(filepath.Dir(originalFileName)),
			"module": cty.StringVal(filepath.Dir(originalFileName)),
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
			Range:    rangeOf(*diag.Subject),
			Source:   "Terraform",
		})
	}

	diags := localDiags(cfg.Locals, originalFileName, variables)
	result = append(result, diags...)

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

	schemata := providerSchema(cfg.ProviderConfigs, originalFileName, variables)
	result = append(result, schemata...)

	schemata = resourceSchema(cfg.ManagedResources, originalFileName, variables)
	result = append(result, schemata...)

	schemata = dataSourceSchema(cfg.DataResources, originalFileName, variables)
	result = append(result, schemata...)

	// spew.Dump(file.ManagedResources[0].Config.Content(nil))

	return result
}

func localDiags(locals []*configs.Local, originalFileName string, variables map[string]cty.Value) []lsp.Diagnostic {
	result := make([]lsp.Diagnostic, 0)

	for _, local := range locals {
		diags := GetLocalsForDiags(*local, filepath.Dir(originalFileName), variables)

		if diags != nil {
			for _, diag := range diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range:    rangeOf(*diag.Subject),
					Source:   "Terraform Schema",
				})
			}
		}
	}

	return result
}

func providerSchema(providers []*configs.Provider, originalFileName string, variables map[string]cty.Value) []lsp.Diagnostic {
	result := make([]lsp.Diagnostic, 0)

	for _, v := range providers {
		providerType := v.Name

		tfSchema := GetProviderSchemaForDiags(providerType, v.Config, filepath.Dir(originalFileName), variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range:    rangeOf(*diag.Subject),
					Source:   "Terraform Schema",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Provider %s does not exist", v.Name),
				Range:    rangeOf(v.NameRange),
				Source:   "Terraform Schema",
			})
		}
	}

	return result
}

func resourceSchema(resources []*configs.Resource, originalFileName string, variables map[string]cty.Value) []lsp.Diagnostic {

	result := make([]lsp.Diagnostic, 0)

	for _, v := range resources {
		resourceType := v.Type

		var providerType string
		if v.ProviderConfigRef != nil {
			providerType = v.ProviderConfigRef.Name
		}

		tfSchema := GetResourceSchemaForDiags(resourceType, v.Config, filepath.Dir(originalFileName), providerType, variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range:    rangeOf(*diag.Subject),
					Source:   "Terraform Schema",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Resource %s does not exist", v.Type),
				Range:    rangeOf(v.TypeRange),
				Source:   "Terraform Schema",
			})
		}
	}

	return result
}

func dataSourceSchema(resources []*configs.Resource, originalFileName string, variables map[string]cty.Value) []lsp.Diagnostic {

	result := make([]lsp.Diagnostic, 0)

	for _, v := range resources {
		resourceType := v.Type
		var providerType string
		if v.ProviderConfigRef != nil {
			providerType = v.ProviderConfigRef.Name
		}

		tfSchema := GetDataSourceSchemaForDiags(resourceType, v.Config, filepath.Dir(originalFileName), providerType, variables)

		if tfSchema != nil {
			for _, diag := range tfSchema.Diags {
				result = append(result, lsp.Diagnostic{
					Severity: lsp.DiagnosticSeverity(diag.Severity),
					Message:  diag.Detail,
					Range:    rangeOf(*diag.Subject),
					Source:   "Terraform",
				})
			}
		} else {
			result = append(result, lsp.Diagnostic{
				Severity: lsp.DiagnosticSeverity(lsp.Error),
				Message:  fmt.Sprintf("Data source %s does not exist", v.Type),
				Range:    rangeOf(v.TypeRange),
				Source:   "Terraform",
			})
		}
	}

	return result
}

func castToV2Range(r oldHCL2.Range) v2.Range {
	return v2.Range{
		Filename: r.Filename,
		Start:    castToV2Pos(r.Start),
		End:      castToV2Pos(r.End),
	}
}

func castToV2Pos(pos oldHCL2.Pos) v2.Pos {
	return v2.Pos{
		Line:   pos.Line,
		Column: pos.Column,
		Byte:   pos.Byte,
	}
}

func rangeOf(r v2.Range) lsp.Range {
	return lsp.Range{
		Start: positionOf(r.Start),
		End:   positionOf(r.End),
	}
}

func positionOf(pos v2.Pos) lsp.Position {
	return lsp.Position{
		Line:      pos.Line - 1,
		Character: pos.Column - 1,
	}
}
