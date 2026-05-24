package main

import (
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift" // импортируем дополнительный анализатор
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
)

func main() {
	multichecker.Main(
		noOsExitMain,
		printf.Analyzer,
		shadow.Analyzer,
		bools.Analyzer,
		tests.Analyzer,
		shift.Analyzer, // добавляем анализатор в вызов multichecker
		structtag.Analyzer,
	)
}
