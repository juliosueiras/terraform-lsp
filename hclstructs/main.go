package hclstructs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"reflect"
)

// Expression hcl
func BinaryOpExpr() reflect.Type {
	return GetType(&hclsyntax.BinaryOpExpr{})
}

func ObjectConsExpr() reflect.Type {
	return GetType(&hclsyntax.ObjectConsExpr{})
}

func FunctionCallExpr() reflect.Type {
	return GetType(&hclsyntax.FunctionCallExpr{})
}

func ScopeTraversalExpr() reflect.Type {
	return GetType(&hclsyntax.ScopeTraversalExpr{})
}

func LiteralValueExpr() reflect.Type {
	return GetType(&hclsyntax.LiteralValueExpr{})
}

func ForExpr() reflect.Type {
	return GetType(&hclsyntax.ForExpr{})
}

func TupleConsExpr() reflect.Type {
	return GetType(&hclsyntax.TupleConsExpr{})
}

func TemplateWrapExpr() reflect.Type {
	return GetType(&hclsyntax.TemplateWrapExpr{})
}

func AnonSymbolExpr() reflect.Type {
	return GetType(&hclsyntax.AnonSymbolExpr{})
}

func SplatExpr() reflect.Type {
	return GetType(&hclsyntax.SplatExpr{})
}

// Traverse hcl
func TraverseAttr() reflect.Type {
	return GetType(hcl.TraverseAttr{})
}

func TraverseIndex() reflect.Type {
	return GetType(hcl.TraverseIndex{})
}

func GetExprStringType(origType reflect.Type) string {

	switch origType {
	// May need recursion
	case BinaryOpExpr():
		return "binary operation"
	case FunctionCallExpr():
		return "function call"
	case ScopeTraversalExpr():
		return "scoped expression"
	case LiteralValueExpr():
		return "literal value"
	case ForExpr():
		return "for loop"
	case TupleConsExpr():
		return "array"
	case TemplateWrapExpr():
		return "string interpolation"
	case ObjectConsExpr():
		return "object"
	case SplatExpr():
		return "splat"
	case AnonSymbolExpr():
		return "anon symbol"
	default:
		return "undefined"
	}

	return "undefined"
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
			for _, newExpr := range expr.ExprList() {
				if newExpr.Range().ContainsPos(posHCL) {
					return newExpr.Variables()
				}
			}
			return nil
		}

	// Need wrapped
	case TemplateWrapExpr():
		expr := expr.(*hclsyntax.TemplateWrapExpr)
		if expr.Range().ContainsPos(posHCL) {
			return expr.Variables()
		}

	case AnonSymbolExpr():
		expr := expr.(*hclsyntax.AnonSymbolExpr)
		if expr.Range().ContainsPos(posHCL) || expr.SrcRange.ContainsPos(posHCL) {
			return expr.Variables()
		}

	case SplatExpr():
		expr := expr.(*hclsyntax.SplatExpr)
		if expr.MarkerRange.ContainsPos(posHCL) || expr.Range().ContainsPos(posHCL) {
			return expr.Source.Variables()
		}

	// Need more check
	case ObjectConsExpr():
		expr := expr.(*hclsyntax.ObjectConsExpr)
		for _, v := range expr.Items {
			if v.KeyExpr.Range().ContainsPos(posHCL) {
			}

			if v.ValueExpr.Range().ContainsPos(posHCL) {
				firstVar := hcl.TraverseAttr{
					Name: v.KeyExpr.(*hclsyntax.ObjectConsKeyExpr).AsTraversal().RootName(),
				}

				vars := hcl.Traversal{
					firstVar,
				}

				origType := reflect.TypeOf(v.ValueExpr)
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

func GetType(t interface{}) reflect.Type {
	return reflect.TypeOf(t)
}
