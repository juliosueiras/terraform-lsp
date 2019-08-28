package tfstructs

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/lang"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/provisioners"
	"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/zclconf/go-cty/cty"
	"path/filepath"
	"strings"
)

type TerraformSchema struct {
	Schema        *providers.Schema
	DecodedSchema cty.Value
	Diags         hcl.Diagnostics
}

type TerraformProvisionerSchema struct {
	Schema        *provisioners.GetSchemaResponse
	DecodedSchema cty.Value
	Diags         hcl.Diagnostics
}

func GetModuleVariables(moduleAddr string, config hcl.Body, targetDir string) (map[string]*configs.Variable, bool) {
	parser := configs.NewParser(nil)

	t, _ := parser.LoadConfigDir(filepath.Join(targetDir, moduleAddr))
	if t == nil || len(t.Variables) == 0 {
		return nil, false
	}

	return t.Variables, true
}

func GetResourceSchema(resourceType string, config hcl.Body, targetDir string, overrideProvider string) *TerraformSchema {
	var provider *Client
	var err error
	if overrideProvider != "" {
		provider, err = GetProvider(overrideProvider, targetDir)
	} else {
		provider, err = GetProvider(resourceType, targetDir)
	}

	if err != nil {
		helper.DumpLog(err)
		return nil
	}

	providerResource, err := provider.GetRawResourceTypeSchema(resourceType)
	if err != nil {
		helper.DumpLog(err)
		provider.Kill()
		return nil
	}

	provider.Kill()

	res2 := providerResource.Block.DecoderSpec()
	// Add Resources and Data Sources & Variables/Functions
	scope := lang.Scope{}

	res, _, diags := hcldec.PartialDecode(config, res2, &hcl.EvalContext{
		// Build Full Tree
		Variables: map[string]cty.Value{
			"path": cty.ObjectVal(map[string]cty.Value{
				"cwd":    cty.StringVal(""),
				"module": cty.StringVal(""),
			}),
			"data":   cty.DynamicVal,
			"var":    cty.DynamicVal, // Need to check for undefined vars
			"module": cty.DynamicVal,
			"local":  cty.DynamicVal,
		},
		Functions: scope.Functions(),
	})

	return &TerraformSchema{
		Schema:        providerResource,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetDataSourceSchema(dataSourceType string, config hcl.Body, targetDir string, overrideProvider string) *TerraformSchema {
	var provider *Client
	var err error
	if overrideProvider != "" {
		provider, err = GetProvider(overrideProvider, targetDir)
	} else {
		provider, err = GetProvider(dataSourceType, targetDir)
	}
	if err != nil {
		helper.DumpLog(err)
		return nil
	}

	providerDataSource, err := provider.GetRawDataSourceTypeSchema(dataSourceType)
	if err != nil {
		helper.DumpLog(err)
		provider.Kill()
		return nil
	}

	provider.Kill()

	res2 := providerDataSource.Block.DecoderSpec()
	scope := lang.Scope{}
	res, _, diags := hcldec.PartialDecode(config, res2, &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"path": cty.ObjectVal(map[string]cty.Value{
				"cwd":    cty.StringVal(""),
				"module": cty.StringVal(""),
			}),
			"data":   cty.DynamicVal,
			"var":    cty.DynamicVal, // Need to check for undefined vars
			"module": cty.DynamicVal,
			"local":  cty.DynamicVal,
		},
		Functions: scope.Functions(),
	})

	return &TerraformSchema{
		Schema:        providerDataSource,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetProvider(providerType string, targetDir string) (*Client, error) {
	if len(strings.Split(providerType, "_")) == 0 {
		return nil, nil
	}
	provider, err := NewClient(strings.Split(providerType, "_")[0], targetDir)
	return provider, err
}

func GetProvisioner(provisionerType string, targetDir string) (*Client, error) {
	provisioner, err := NewProvisionerClient(provisionerType, targetDir)
	return provisioner, err
}

func GetProvisionerSchema(provisionerType string, config hcl.Body, targetDir string) *TerraformProvisionerSchema {
	provisioner, err := GetProvider(provisionerType, targetDir)
	if err != nil {
		helper.DumpLog(err)
		return nil
	}

	provisionerSchema := provisioner.provisioner.GetSchema()

	provisioner.Kill()

	res2 := provisionerSchema.Provisioner.DecoderSpec()
	scope := lang.Scope{}

	res, _, diags := hcldec.PartialDecode(config, res2, &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"data":   cty.DynamicVal,
			"var":    cty.DynamicVal, // Need to check for undefined vars
			"module": cty.DynamicVal,
		},
		Functions: scope.Functions(),
	})

	return &TerraformProvisionerSchema{
		Schema:        &provisionerSchema,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetProviderSchema(providerType string, config hcl.Body, targetDir string) *TerraformSchema {
	provider, err := GetProvider(providerType, targetDir)
	if err != nil {
		helper.DumpLog(err)
		return nil
	}

	providerSchema := provider.provider.GetSchema().Provider

	provider.Kill()

	res2 := providerSchema.Block.DecoderSpec()
	scope := lang.Scope{}

	res, _, diags := hcldec.PartialDecode(config, res2, &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"path": cty.ObjectVal(map[string]cty.Value{
				"cwd":    cty.StringVal(""),
				"module": cty.StringVal(""),
			}),
			"data":   cty.DynamicVal,
			"var":    cty.DynamicVal, // Need to check for undefined vars
			"module": cty.DynamicVal,
			"local":  cty.DynamicVal,
		},
		Functions: scope.Functions(),
	})

	return &TerraformSchema{
		Schema:        &providerSchema,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetAllConfigs(filePath string, tempFilePath string) *configs.Module {
	parser := configs.NewParser(nil)
	fileURL := strings.Replace(filePath, "file://", "", 1)

	fileDir := filepath.Dir(fileURL)
	res, _ := filepath.Glob(fileDir + "/*.tf")
	var resultFiles []*configs.File

	for _, v := range res {
		if fileURL == v {
			continue
		}

		cFile, _ := parser.LoadConfigFile(v)

		resultFiles = append(resultFiles, cFile)
	}

	tempConfig, _ := parser.LoadConfigFile(tempFilePath)
	resultFiles = append(resultFiles, tempConfig)

	files, _ := configs.NewModule(resultFiles, nil)

	return files
}
