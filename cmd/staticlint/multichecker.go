// Mycheck предназначен для анализа кода, используя дерево AST и статистические анализаторы
// для запуска необходимо скомпилировать пакет и передать бинарнику путь, который будем проверять
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift" // импортируем дополнительный анализатор
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"honnef.co/go/tools/staticcheck"
)

var noOsExitMain = &analysis.Analyzer{
	Name: "noosexitinmain",
	Doc:  "запрещает прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func main() {
	// Получаем все анализаторы класса SA
	analyzersstatic := staticcheck.Analyzers

	// Передаем их в multichecker
	var checks = make([]*analysis.Analyzer, 0)
	for _, a := range analyzersstatic {
		checks = append(checks, a.Analyzer)
	}
	// добавляем в multichecker анализаторы пакета go/analysis
	checks = append(checks, printf.Analyzer,
		shadow.Analyzer,
		bools.Analyzer,
		tests.Analyzer,
		shift.Analyzer,
		structtag.Analyzer)
	// добавлеяем в multichecker свой анализаторо проверки os.Exit in main
	checks = append(checks, noOsExitMain)

	// запускаем все анализаторы
	multichecker.Main(checks...)

}

func run(pass *analysis.Pass) (interface{}, error) {
	// Проверяем, что анализируем пакет "main"
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			// Ищем объявление функции
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Проверяем, что функция называется "main"
			if fn.Name.Name != "main" || fn.Body == nil {
				return true
			}

			// Проходим по телу функции и ищем вызовы os.Exit
			for _, stmt := range fn.Body.List {
				if isOsExitCall(stmt, pass) {
					pass.Reportf(stmt.Pos(), "прямой вызов os.Exit() в функции main пакета main запрещен")
				}
			}
			return false
		})
	}
	return nil, nil
}

// Проверяет, является ли выражение вызовом os.Exit ... ,pass *analysis.Pass
func isOsExitCall(stmt ast.Stmt, pass *analysis.Pass) bool {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	b := pass.TypesInfo.Uses[ident]
	return b.Name() == "os" && selExpr.Sel.Name == "Exit"
}
