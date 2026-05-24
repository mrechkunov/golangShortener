package main

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Pass struct {
	Fset         *token.FileSet // информация о позиции токенов
	Files        []*ast.File    // AST для каждого файла
	OtherFiles   []string       // имена файлов не на Go в пакете
	IgnoredFiles []string       // имена игнорируемых исходных файлов в пакете
	Pkg          *types.Package // информация о типах пакета
	TypesInfo    *types.Info    // информация о типах в AST
}

var noOsExitMain = &analysis.Analyzer{
	Name: "noosexitmain",
	Doc:  "check os.Exit run at main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	// Анализируем только пакет "main"
	if pass.Pkg == nil || pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем объявление функции
			fnDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Проверяем, что функция называется "main"
			if fnDecl.Name.Name != "main" || fnDecl.Body == nil {
				return true
			}

			// Проходим по всем выражениям внутри тела функции main
			for _, stmt := range fnDecl.Body.List {
				exprStmt, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}
				callExpr, ok := exprStmt.X.(*ast.CallExpr)
				if !ok {
					continue
				}

				// Проверяем, что вызывается именно функция os.Exit
				if isOsExit(pass, callExpr) {
					pass.Reportf(callExpr.Pos(), "прямой вызов os.Exit в функции main пакета main запрещен")
				}
			}
			return true
		})
	}

	return nil, nil
}

// isOsExit проверяет, что вызываемая функция — это os.Exit
func isOsExit(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Проверяем имя вызываемой функции (Exit)
	if sel.Sel == nil || sel.Sel.Name != "Exit" {
		return false
	}

	// Проверяем пакет (os)
	if ident, ok := sel.X.(*ast.Ident); ok {
		obj := pass.TypesInfo.ObjectOf(ident)
		if pkgName, ok := obj.(*types.PkgName); ok {
			if pkgName.Imported() != nil && pkgName.Imported().Path() == "os" {
				return true
			}
		}
	}

	return false
}
