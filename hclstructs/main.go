package hclstructs

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/juliosueiras/terraform-lsp/helper"
	"reflect"
)

func BinaryOpExpr() reflect.Type {
	return helper.GetType(&hclsyntax.BinaryOpExpr{})
}

func ObjectConsExpr() reflect.Type {
	return helper.GetType(&hclsyntax.ObjectConsExpr{})
}

func FunctionCallExpr() reflect.Type {
	return helper.GetType(&hclsyntax.FunctionCallExpr{})
}

func ScopeTraversalExpr() reflect.Type {
	return helper.GetType(&hclsyntax.ScopeTraversalExpr{})
}

func LiteralValueExpr() reflect.Type {
	return helper.GetType(&hclsyntax.LiteralValueExpr{})
}

func ForExpr() reflect.Type {
	return helper.GetType(&hclsyntax.ForExpr{})
}

func TupleConsExpr() reflect.Type {
	return helper.GetType(&hclsyntax.TupleConsExpr{})
}

func TemplateWrapExpr() reflect.Type {
	return helper.GetType(&hclsyntax.TemplateWrapExpr{})
}

func GetExprVariables(origType reflect.Type, expr hcl.Expression, posHCL hcl.Pos) []hcl.Traversal {

	switch origType {
	// May need recursion
	case BinaryOpExpr():
		expr := expr.(*hclsyntax.BinaryOpExpr)
		if expr.LHS.Range().ContainsPos(posHCL) {
			return expr.LHS.Variables()
		} else if expr.RHS.Range().ContainsPos(posHCL) {
			return expr.RHS.Variables()
		}

	// Need more check
	case FunctionCallExpr():
		expr := expr.(*hclsyntax.FunctionCallExpr)
		for _, arg := range expr.ExprCall().Arguments {
			if arg.Range().ContainsPos(posHCL) {
				return arg.Variables()
			}
		}
	case ScopeTraversalExpr():
		expr := expr.(*hclsyntax.ScopeTraversalExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}
	case LiteralValueExpr():
		expr := expr.(*hclsyntax.LiteralValueExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}

	// Need more check
	case ForExpr():
		expr := expr.(*hclsyntax.ForExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}

	// Need more check
	case TupleConsExpr():
		expr := expr.(*hclsyntax.TupleConsExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}

	// Need wrapped
	case TemplateWrapExpr():
		expr := expr.(*hclsyntax.TemplateWrapExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}

	// Need more check
	case ObjectConsExpr():
		expr := expr.(*hclsyntax.ObjectConsExpr)
		for _, v := range expr.Items {
			if v.KeyExpr.Range().ContainsPos(posHCL) {
				helper.DumpLog(v)
			}

			if v.ValueExpr.Range().ContainsPos(posHCL) {
				firstVar := hcl.TraverseAttr{
					Name: v.KeyExpr.(*hclsyntax.ObjectConsKeyExpr).AsTraversal().RootName(),
				}
				vars := hcl.Traversal{
					firstVar,
				}

				origType := reflect.TypeOf(expr)
				helper.DumpLog(origType)
				res := GetExprVariables(origType, v.ValueExpr, posHCL)
				resultAttrs := vars
				if len(res) != 0 && res != nil {
					resultAttrs = append(vars, res[0]...)
				}
				return []hcl.Traversal{
					resultAttrs,
				}
			}
		}

	default:
		return nil
	}

	return nil
}
