package tfstructs

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/providers"
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

func GetModuleVariables(moduleAddr string, config hcl.Body, targetDir string) (map[string]*configs.Variable, bool) {
	parser := configs.NewParser(nil)

	t, _ := parser.LoadConfigDir(filepath.Join(targetDir, moduleAddr))
	if t == nil || len(t.Variables) == 0 {
		return nil, false
	}

	return t.Variables, true
}

func GetResourceSchema(resourceType string, config hcl.Body, targetDir string) *TerraformSchema {
	provider, err := GetProvider(resourceType, targetDir)
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
	res, diags := hcldec.Decode(config, res2, nil)

	return &TerraformSchema{
		Schema:        providerResource,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetDataSourceSchema(dataSourceType string, config hcl.Body, targetDir string) *TerraformSchema {
	provider, err := GetProvider(dataSourceType, targetDir)
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
	res, diags := hcldec.Decode(config, res2, nil)

	return &TerraformSchema{
		Schema:        providerDataSource,
		DecodedSchema: res,
		Diags:         diags,
	}
}

func GetProvider(resourceType string, targetDir string) (*Client, error) {
	if len(strings.Split(resourceType, "_")) == 0 {
		return nil, nil
	}
	provider, err := NewClient(strings.Split(resourceType, "_")[0], targetDir)
	return provider, err
}
