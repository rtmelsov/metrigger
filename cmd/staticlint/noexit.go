package staticlint

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var NoExistAnalyzer = &analysis.Analyzer{
	Name: "noexit",
	Doc:  "Can't use Exit function in main package",
	Run:  noExitRun,
}

func noExitRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		f := pass.Fset.File(file.Pos())
		if !strings.HasSuffix(f.Name(), "main.go") {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			pkgIdent, ok := sel.X.(*ast.Ident)
			if !ok || pkgIdent.Name != "os" || sel.Sel.Name != "Exit" {
				return true
			}

			pass.Reportf(sel.Pos(), "direct os.Exit in main is forbidden")
			return false
		})
	}
	return nil, nil
}
