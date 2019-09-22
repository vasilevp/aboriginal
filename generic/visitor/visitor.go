package visitor

import (
	"go/ast"
	"go/token"
	"log"
	"regexp"

	"github.com/vasilevp/aboriginal/generic/gtype"
)

type Visitor struct {
	parent         *Visitor
	parentNode     ast.Node
	genericTypes   map[string]gtype.Type
	generatedDecls map[string]ast.Decl
}

func New(a *ast.File) Visitor {
	return Visitor{
		parentNode:     a,
		genericTypes:   map[string]gtype.Type{},
		generatedDecls: map[string]ast.Decl{},
	}
}

var templateRegex = regexp.MustCompile(`([\p{L}0-9]+)ᐸ([\p{L}0-9]+)ᐳ`)

func (p Visitor) createGenericType(n *ast.TypeSpec) bool {
	switch t := n.Type.(type) {
	case *ast.StructType:
		matches := templateRegex.FindStringSubmatch(n.Name.Name)
		template, typename := matches[1], matches[2]
		log.Printf("created generic type %q\n", template)
		p.genericTypes[template] = gtype.Type{
			Args:    []string{typename},
			Node:    t,
			Methods: map[string]*ast.FuncDecl{},
		}

		return true // don't process contents of struct
	}

	return false
}

func renameField(f *ast.Field, name string) *ast.Field {
	return &ast.Field{
		Doc:   f.Doc,
		Names: f.Names,
		Type: &ast.Ident{
			Name: name,
		},
		Tag: f.Tag,
		// Comment: f.Comment,
	}
}

func createTypeDecl(name string, fields *ast.FieldList) ast.Decl {
	return &ast.GenDecl{
		Doc: nil,
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: name,
				},
				Type: &ast.StructType{
					Fields: fields,
				},
			},
		},
	}
}

func renameFuncReceiver(f *ast.FuncDecl, name string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc:  f.Doc,
		Name: f.Name,
		Type: f.Type,
		Body: f.Body,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				renameField(f.Recv.List[0], name),
			},
		},
	}
}

func (p Visitor) accessGenericType(n *ast.Ident) {
	matches := templateRegex.FindStringSubmatch(n.Name)
	generic, typename := matches[1], matches[2]
	g, ok := p.genericTypes[generic]
	if !ok {
		log.Printf("%q is not a generic type\n", generic)
		return
	}

	log.Printf("accessed generic type %q<%q>", generic, typename)

	if _, ok := p.generatedDecls[n.Name]; ok {
		log.Printf("type %q is already generated, bailing out", n.Name)
		return
	}

	f := &ast.FieldList{}

	for _, v := range g.Node.Fields.List {
		id, ok := v.Type.(*ast.Ident)
		if !ok {
			f.List = append(f.List, v)
			continue
		}

		name := id.Name
		if name == g.Args[0] {
			log.Printf("substituting %s with %s", name, typename)
			name = typename
		}

		f.List = append(f.List, renameField(v, name))
	}

	decl := createTypeDecl(n.Name, f)
	p.generatedDecls[n.Name] = decl

	for _, v := range g.Methods {
		p.generatedDecls[n.Name+"."+v.Name.Name] = renameFuncReceiver(v, n.Name)
	}
}

func (p *Visitor) generateType() {

}

func (p *Visitor) createGenericMethod(f *ast.FuncDecl) {
	name := getReceiverType(f).Name
	matches := templateRegex.FindStringSubmatch(name)
	template := matches[1]
	p.genericTypes[template].Methods[f.Name.Name] = f
}

func getReceiverType(f *ast.FuncDecl) *ast.Ident {
	if f.Recv == nil || len(f.Recv.List) == 0 {
		return nil
	}

	t, ok := f.Recv.List[0].Type.(*ast.Ident)
	if !ok {
		tt, ok := f.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}

		t, ok = tt.X.(*ast.Ident)
		if !ok {
			return nil
		}
	}

	return t
}

func (p Visitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl:
		t := getReceiverType(n)
		if t == nil {
			break
		}

		if !templateRegex.MatchString(t.Name) {
			break
		}

		p.createGenericMethod(n)

		return nil

	case *ast.Ident:
		if !templateRegex.MatchString(n.Name) {
			break
		}

		p.accessGenericType(n)

	case *ast.TypeSpec:
		if !templateRegex.MatchString(n.Name.Name) {
			break
		}

		if p.createGenericType(n) {
			return nil
		}
	}

	result := Visitor{
		parent:         &p,
		parentNode:     n,
		genericTypes:   p.genericTypes,
		generatedDecls: p.generatedDecls,
	}

	return result
}

func (p Visitor) GeneratedDecls() map[string]ast.Decl {
	return p.generatedDecls
}
