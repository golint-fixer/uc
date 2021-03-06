package sem

import (
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewmew/uc/ast"
	"github.com/mewmew/uc/ast/astutil"
	"github.com/mewmew/uc/sem/errors"
	"github.com/mewmew/uc/types"
)

// universePos specifies a pseudo-position used for identifiers declared in the
// universe scope.
const universePos = -1

// resolve performs identifier resolution, mapping identifiers to corresponding
// declarations.
func resolve(file *ast.File, scopes map[ast.Node]*Scope) error {
	// TODO: Verify that type keywords cannot be redeclared.

	// Pre-pass, add keyword types and universe scope.
	universe := NewScope(nil)
	charIdent := &ast.Ident{NamePos: universePos, Name: "char"}
	charDecl := &ast.TypeDef{DeclType: charIdent, TypeName: charIdent, Val: &types.Basic{Kind: types.Char}}
	charIdent.Decl = charDecl
	intIdent := &ast.Ident{NamePos: universePos, Name: "int"}
	intDecl := &ast.TypeDef{DeclType: intIdent, TypeName: intIdent, Val: &types.Basic{Kind: types.Int}}
	intIdent.Decl = intDecl
	voidIdent := &ast.Ident{NamePos: universePos, Name: "void"}
	voidDecl := &ast.TypeDef{DeclType: voidIdent, TypeName: voidIdent, Val: &types.Basic{Kind: types.Void}}
	voidIdent.Decl = voidDecl
	universeDecls := []*ast.TypeDef{
		charDecl,
		intDecl,
		voidDecl,
	}
	for _, decl := range universeDecls {
		if err := universe.Insert(decl); err != nil {
			return errutil.Err(err)
		}
	}

	// First pass, add global declarations to file scope.
	fileScope := NewScope(universe)
	scopes[file] = fileScope
	fileScope.IsDef = func(decl ast.Decl) bool {
		// Consider variable declarations as tentative definitions; i.e. return
		// false, unless variable definition.
		return decl.Value() != nil
	}
	for _, decl := range file.Decls {
		if err := fileScope.Insert(decl); err != nil {
			return errutil.Err(err)
		}
	}

	// skip specifies that the block statement body of a function declaration
	// should skip creating a nested scope, as it has already been created by its
	// function declaration, so that function parameters are placed within the
	// correct scope.
	skip := false

	// scope specifies the current lexical scope.
	scope := fileScope

	// resolve performs identifier resolution, mapping identifiers to the
	// corresponding declarations of the closest lexical scope.
	resolve := func(n ast.Node) error {
		switch n := n.(type) {
		case ast.Decl:
			// Insert declaration into the scope if not already added by the
			// file scope pre-pass.
			if scope != fileScope {
				if err := scope.Insert(n); err != nil {
					return errutil.Err(err)
				}
			}
			// Create nested scope for function definitions.
			if fn, ok := n.(*ast.FuncDecl); ok {
				if astutil.IsDef(fn) {
					skip = true
				}
				scope = NewScope(scope)
				scopes[fn] = scope
			}
		case *ast.BlockStmt:
			if !skip {
				scope = NewScope(scope)
				scopes[n] = scope
			}
			skip = false
		case *ast.Ident:
			decl, ok := scope.Lookup(n.Name)
			if !ok {
				return errors.Newf(n.Start(), "undeclared identifier %q", n)
			}
			n.Decl = decl
		}
		return nil
	}

	// after reverts to the outer scope after traversing block statements.
	after := func(n ast.Node) error {
		if _, ok := n.(*ast.BlockStmt); ok {
			scope = scope.Outer
		} else if fn, ok := n.(*ast.FuncDecl); ok && !astutil.IsDef(fn) {
			scope = scope.Outer
		}
		return nil
	}

	// Walk the AST of the given file to resolve identifiers.
	if err := astutil.WalkBeforeAfter(file, resolve, after); err != nil {
		return errutil.Err(err)
	}

	return nil
}
