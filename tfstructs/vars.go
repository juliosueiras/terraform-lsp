package tfstructs

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/juliosueiras/terraform-lsp/helper"
	"github.com/sourcegraph/go-lsp"
)

type GetVarAttributeRequest struct {
	Variables hcl.Traversal
	Result    []lsp.CompletionItem
	Files     *configs.Module
	Config    hcl.Body
	FileDir   string
}

func GetVarAttributeCompletion(request GetVarAttributeRequest) []lsp.CompletionItem {
	helper.DumpLog("hi")
	helper.DumpLog(request.Variables)
	if request.Variables.RootName() == "var" {
		vars := request.Variables

		request.Result = helper.ParseVariables(vars[1:], request.Files.Variables, request.Result)
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
				request.Result = append(request.Result, lsp.CompletionItem{
					Label:  v.Type,
					Detail: " data resource",
				})
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
