package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Tree map[string]interface{}

func newTree() Tree {
	return make(map[string]interface{})
}

func getRealTypeName(expr ast.Expr) string {

	switch expr.(type) {
	case *ast.StarExpr:
		return getRealTypeName(expr.(*ast.StarExpr).X)

		//return nil
	case *ast.Ident:
		return expr.(*ast.Ident).Name
	default:
		return "UnknownReceiver"
	}

	//return nil
}

// func getFuncParameterNameList(funcElement *ast.FuncDecl) {
// 	params := make([]string)
// 	for i := 0; i < funcNode.Type.Params.NumFields(); i++ {
// 		params = append(params, funcNode.Type.Params.List[0].Names[0])
// 	}
// }

func (t Tree) Visit(node ast.Node) (w ast.Visitor) {

	switch node.(type) {
	case *ast.FuncDecl:
		funcNode := node.(*ast.FuncDecl)
		funcTree := newTree()
		funcTree["line"] = funcNode.Pos()
		if funcNode.Recv.NumFields() > 0 {
			funcTree["receiver"] = getRealTypeName(funcNode.Recv.List[0].Type)
		}
		funcTree["elemtype"] = "func"
		funcTree["name"] = funcNode.Name.Name
		funcTree["public"] = funcNode.Name.IsExported()
		t[funcNode.Name.Name] = funcTree

	case *ast.TypeSpec:
		typeNode := node.(*ast.TypeSpec)
		typeTree := newTree()

		typeTree["name"] = typeNode.Name.Name
		typeTree["elemtype"] = "type"
		typeTree["public"] = typeNode.Name.IsExported()

		t[typeNode.Name.Name] = typeTree
	}
	return t
}

func parseFile(inputFile string) int {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)

	outputTree := newTree()
	outputTree["file"] = inputFile

	ast.Walk(&outputTree, tree)
	if err == nil {
		outputTree["package"] = tree.Name.Name
		enc := json.NewEncoder(os.Stdout)
		fmt.Println(enc.Encode(outputTree))
	}

	return 0
}

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
