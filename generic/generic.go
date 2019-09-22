package generic

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"

	"github.com/vasilevp/aboriginal/generic/visitor"

	"github.com/pkg/errors"
)

func copyImports(src, dest *ast.File) {
	for _, v := range src.Decls {
		v, ok := v.(*ast.GenDecl)
		if !ok {
			break
		}

		for _, spec := range v.Specs {
			if _, ok := spec.(*ast.ImportSpec); !ok {
				return
			}
		}

		dest.Decls = append(dest.Decls, v)
	}
}

func Process(in interface{}, out io.Writer, filename string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, in, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "cannot parse source")
	}

	v := visitor.New(f)
	ast.Walk(v, f)

	newf := &ast.File{
		Doc:        f.Doc,
		Name:       f.Name,
		Decls:      nil,
		Scope:      f.Scope,
		Imports:    f.Imports,
		Unresolved: f.Unresolved,
		// Comments:   f.Comments,
	}

	for _, v := range f.Unresolved {
		log.Printf("U %v", v.Name)
	}

	copyImports(f, newf)

	for _, t := range v.GeneratedDecls() {
		newf.Decls = append(newf.Decls, t)
	}

	err = format.Node(out, fset, newf)
	if err != nil {
		return errors.Wrap(err, "cannot format source")
	}

	return nil
}
