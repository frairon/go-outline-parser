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

func (t Tree) Visit(node ast.Node) (w ast.Visitor) {
	fmt.Printf("%T\n", node)
	switch node.(type) {
	case *ast.FuncDecl:
		funcNode := node.(*ast.FuncDecl)
		funcTree := newTree()
		funcTree["line"] = funcNode.Pos()

		t[funcNode.Name.Name] = funcTree
	//t["package"] = node.(ast.Package).Name
	default:
		fmt.Println("unknown node")
	}
	fmt.Println(node)
	return t
}

func parseFile(inputFile string) int {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)

	fmt.Println(tree)
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
