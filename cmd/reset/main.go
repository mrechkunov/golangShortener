package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

func getNodeSource(fset *token.FileSet, node ast.Node) (string, error) {
	// Получаем координаты начала и конца узла в байтах
	pos := fset.Position(node.Pos())
	end := fset.Position(node.End())

	// Читаем весь файл в память
	content, err := os.ReadFile(pos.Filename)
	if err != nil {
		return "", err
	}

	// Вырезаем точный кусок кода по байтовым смещениям
	return string(content[pos.Offset:end.Offset]), nil
}

// NodeToString превращает любой узел AST в красивый Go-код
func NodeToString(fset *token.FileSet, node ast.Node) (string, error) {
	var buf bytes.Buffer
	// Node() принимает fset, сам узел и пишет результат в буфер
	err := format.Node(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func main() {
	// Конфигурация: загружаем AST и базовую информацию о типах
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes,
		Tests: true, // Включаем тесты (можно отключить для production)
	}

	// Загружаем пакеты для текущей директории ("./...")
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatalf("Ошибка загрузки пакетов: %v", err)
	}

	// Обходим все загруженные пакеты
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		// Каждый файл в пакете имеет свое собственное AST (*ast.File)
		for _, file := range pkg.Syntax {
			for _, gr := range file.Comments {
				for _, c := range gr.List {
					if strings.HasPrefix(c.Text, "//generate:reset") {
						fmt.Println(pkg.Name)
						fmt.Println(pkg.Fset.Position(c.Slash).String(), c.Text)

					}
				}
			}
		}
	})
}
