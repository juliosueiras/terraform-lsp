package hclstructs

import (
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/juliosueiras/terraform-lsp/helper"
	"reflect"
)

func BinaryOpExpr() reflect.Type {
	return helper.GetType(&hclsyntax.BinaryOpExpr{})
}

func FunctionCallExpr() reflect.Type {
	return helper.GetType(&hclsyntax.FunctionCallExpr{})
}
