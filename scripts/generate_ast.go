package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defineAst("./cmd/myinterpreter/ast", "Expr", []string{
		"Assign   : Value Expr, Name Token",
		"Binary   : Left Expr, Right Expr, Operator Token",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Unary    : Right Expr, Operator Token",
		"Variable : Name Token",
	},
	)
	defineAst("./cmd/myinterpreter/ast", "Stmt", []string{
		"Expression : Expression Expr",
		"Print      : Expression Expr",
		"Var        : Initializer Expr, Name Token",
	},
	)
}

func defineAst(outputDir string, baseName string, types []string) {
	path := filepath.Join(outputDir, strings.ToLower(baseName)+".go")
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("package ast\n\n")

	defineStruct(file, baseName)

	defineVisitor(file, baseName, types)

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		fields := strings.Split(strings.TrimSpace(parts[1]), ",")

		defineType(file, baseName, className, fields)
	}
}

func defineType(file *os.File, baseName string, className string, fields []string) {
	fmt.Fprintf(file, "type %s struct {\n", className)

	for _, field := range fields {
		fieldParts := strings.Fields(strings.TrimSpace(field))

		if len(fieldParts) == 2 {
			fmt.Fprintf(file, "\t%s %s\n", fieldParts[0], fieldParts[1])
		}
	}

	file.WriteString("}\n\n")

	alias := strings.ToLower(string(className[0]))
	fmt.Fprintf(file, "func (%s *%s) Accept(visitor %sVisitor) interface{} {\n", alias, className, baseName)
	fmt.Fprintf(file, "\treturn visitor.Visit%s%s(%s)\n", className, baseName, alias)

	file.WriteString("}\n\n")
}

func defineVisitor(file *os.File, baseName string, types []string) {
	fmt.Fprintf(file, "type %sVisitor interface {\n", baseName)

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])

		fmt.Fprintf(file, "\tVisit%s%s(expt *%s) interface{}\n", className, baseName, className)
	}

	file.WriteString("}\n\n")
}

func defineStruct(file *os.File, baseName string) {
	fmt.Fprintf(file, "type %s interface {\n", baseName)
	fmt.Fprintf(file, "\tAccept(visitor %sVisitor) interface{}\n", baseName)
	file.WriteString("}\n\n")
}
