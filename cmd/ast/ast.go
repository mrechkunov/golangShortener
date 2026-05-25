package asttest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func asttest() {
	// Пустой набор файлов
	fset := token.NewFileSet()
	// Создаем AST и передаем FileSet, чтобы запомнить «маппинг» токенов.

	// В качестве src можно использовать string, []byte или io.Reader
	f, err := parser.ParseFile(fset, "/Users/michaelrechkunov/develop/yaPracticum/advancedGolangDeveloper/golangShortener/cmd/shortener/main.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		// проверяем, какой конкретный тип лежит в узле
		switch x := n.(type) {
		case *ast.CallExpr:
			// ast.CallExpr представляет вызов функции или метода
			fmt.Printf("CallExpr %v: ", fset.Position(x.Fun.Pos()))
			printer.Fprint(os.Stdout, fset, x)
			fmt.Println()
		case *ast.FuncDecl:
			// ast.FuncDecl представляет декларацию функции
			fmt.Printf("FuncDecl %s %v: ", x.Name.Name, fset.Position(x.Pos()))
			printer.Fprint(os.Stdout, fset, x)
			fmt.Println()
		}
		return true
	})
	// if err = ast.Print(fset, f); err != nil {
	// 	panic(err)
	// }
}
