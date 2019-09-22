package gtype

import "go/ast"

type Type struct {
	Args    []string
	Methods map[string]*ast.FuncDecl
	Node    *ast.StructType
	Parent  ast.Node
}
