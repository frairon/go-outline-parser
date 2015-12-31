package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Entry map[string]interface{}

type FileOutline struct {
	FileSet     *token.FileSet `json:"-"`
	Filename    string
	Packagename string

	Entries map[string]Entry
}

func newEntry() Entry {
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

}

func (o FileOutline) Visit(node ast.Node) (w ast.Visitor) {

	switch node.(type) {
	case *ast.FuncDecl:
		funcNode := node.(*ast.FuncDecl)
		funcTree := newEntry()

		o.setPosition(funcTree, funcNode.Pos())

		if funcNode.Recv.NumFields() > 0 {
			funcTree["Receiver"] = getRealTypeName(funcNode.Recv.List[0].Type)
		}
		funcTree["Elemtype"] = "func"
		funcTree["Name"] = funcNode.Name.Name
		funcTree["Public"] = funcNode.Name.IsExported()
		o.Entries[funcNode.Name.Name] = funcTree

	case *ast.TypeSpec:
		typeNode := node.(*ast.TypeSpec)
		typeTree := newEntry()

		typeTree["Name"] = typeNode.Name.Name
		typeTree["Elemtype"] = "type"
		typeTree["Public"] = typeNode.Name.IsExported()

		o.setPosition(typeTree, typeNode.Pos())

		o.Entries[typeNode.Name.Name] = typeTree
	}
	return o
}

func (o *FileOutline) setPosition(entry Entry, pos token.Pos) {
	entry["Line"] = 0
	entry["Column"] = 0

	if pos.IsValid() {
		fpos := o.FileSet.Position(pos)
		if fpos.IsValid() {
			entry["Line"] = fpos.Line
			entry["Column"] = fpos.Column
		}
	}
}

func parseFile(inputFile string) int {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)
	outline := FileOutline{
		Entries: make(map[string]Entry),
		FileSet: fset,
	}

	outline.Filename = inputFile

	ast.Walk(&outline, tree)
	if err == nil {
		outline.Packagename = tree.Name.Name
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(outline)
		return 0
	}

	fmt.Println("Error parsing go file.")
	return 1

}
