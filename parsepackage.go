package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type PackageDef string

func (pd *PackageDef) Visit(node ast.Node) (w ast.Visitor) {
	switch node.(type) {
	case *ast.Package:
		*pd = node.(*ast.Package).Name
	}
	return pd
}

func parsePackage(inputFile string) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.PackageClauseOnly)
	var pkgDef PackageDef

	ast.Walk(&pkgDef, tree)

	fmt.Prinln(pkgDef)
}
