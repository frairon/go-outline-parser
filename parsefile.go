package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Entry struct {
	Elemtype string
	Receiver string `json:",omitempty"`
	Name     string
	Public   bool
	Line     int `json:",omitempty"`
	Column   int `json:",omitempty"`
}

type FileOutline struct {
	FileSet     *token.FileSet `json:"-"`
	Filename    string
	Packagename string

	Entries []*Entry
}

func newEntry() *Entry {
	return new(Entry)
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

func (o *FileOutline) Visit(node ast.Node) (w ast.Visitor) {
	switch typedNode := node.(type) {
	case *ast.FuncDecl:
		funcTree := newEntry()

		o.setPosition(funcTree, typedNode.Name.Pos())

		if typedNode.Recv.NumFields() > 0 {
			funcTree.Receiver = getRealTypeName(typedNode.Recv.List[0].Type)
		}
		funcTree.Elemtype = "func"
		funcTree.Name = typedNode.Name.Name
		funcTree.Public = typedNode.Name.IsExported()
		o.Entries = append(o.Entries, funcTree)

		// return nil so it does not recurse into functions
		return nil
	// case *ast.InterfaceType:
	// 	fmt.Printf("%+v", typedNode.Methods.List[0].Names)
	// 	Interface
	case *ast.TypeSpec:
		typeTree := newEntry()
		typeTree.Name = typedNode.Name.Name
		typeTree.Elemtype = "type"
		typeTree.Public = typedNode.Name.IsExported()

		o.setPosition(typeTree, typedNode.Pos())

		switch typedType := typedNode.Type.(type) {
		case *ast.InterfaceType:
			typeTree.Elemtype = "interface"
			o.createMethodsForInterface(typedType, typedNode.Name.Name)
		}

		o.Entries = append(o.Entries, typeTree)

	case *ast.ValueSpec:
		for _, name := range typedNode.Names {
			valueTree := newEntry()
			valueTree.Elemtype = "variable"
			valueTree.Name = name.Name
			o.setPosition(valueTree, name.Pos())
			o.Entries = append(o.Entries, valueTree)
		}
	}
	return o
}

func (o *FileOutline) createMethodsForInterface(iface *ast.InterfaceType, receiver string) {
	for _, field := range iface.Methods.List {
		if len(field.Names) < 1 {
			continue
		}
		ifaceFunc := newEntry()
		ifaceFunc.Name = field.Names[0].Name
		ifaceFunc.Elemtype = "func"
		ifaceFunc.Receiver = receiver
		o.setPosition(ifaceFunc, field.Type.Pos())

		o.Entries = append(o.Entries, ifaceFunc)
	}
}

func (o *FileOutline) setPosition(entry *Entry, pos token.Pos) {

	if pos.IsValid() {
		fpos := o.FileSet.Position(pos)
		if fpos.IsValid() {
			entry.Line = fpos.Line
			entry.Column = fpos.Column
		}
	}
}

func parseFile(inputFile string) (string, error) {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, inputFile, nil, parser.AllErrors)
	if err != nil {
		return "", fmt.Errorf("Error parsing the file: %v", err)
	}
	outline := FileOutline{
		FileSet: fset,
	}

	outline.Filename = inputFile

	ast.Walk(&outline, tree)

	outline.Packagename = tree.Name.Name
	output := bytes.Buffer{}
	enc := json.NewEncoder(&output)
	err = enc.Encode(outline)
	if err != nil {
		return "", fmt.Errorf("Error encoding the symbol tree: %v", err)
	}
	return output.String(), nil
}
