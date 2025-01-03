package xast

import "go/ast"

func IsTemplComponent(results *ast.FieldList) bool {
	if results == nil || len(results.List) != 1 {
		return false
	}

	// check if return type is templ.Component
	selectorExpr, ok := results.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "templ" && selectorExpr.Sel.Name == "Component"
}
