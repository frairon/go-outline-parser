package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type PackageDef string

func (pd *PackageDef) Visit(node ast.Node) (w ast.Visitor) {
	fmt.Printf("found node %T\n", node)
	switch node.(type) {
	case *ast.Package:
		var name string = node.(*ast.Package).Name
		fmt.Println("found package", name)
	}
	return pd
}

func parsePackage(inputFile string) int {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.PackageClauseOnly)

	if err != nil {
		return 1
	}
	var pkgDef PackageDef

	ast.Walk(&pkgDef, tree)

	fmt.Println(pkgDef)

	return 0
}
