package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	types := []string{
		"Binary   : Left Expr, Right Expr, Operator Token",
		"Grouping : Expression Expr",
		"Literal  : Value interface{}",
		"Unary    : Right Expr, Operator Token",
	}

	defineAst("./cmd/myinterpreter/ast", "classes", types)
}

func defineAst(outputDir string, baseName string, types []string) {
	path := filepath.Join(outputDir, baseName+".go")
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("package ast\n\n")

	defineVisitor(file, types)

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		fields := strings.Split(strings.TrimSpace(parts[1]), ",")

		defineType(file, className, fields)
	}
}

func defineType(file *os.File, className string, fields []string) {
	fmt.Fprintf(file, "type %s struct {\n", className)

	for _, field := range fields {
		fieldParts := strings.Fields(strings.TrimSpace(field))

		if len(fieldParts) == 2 {
			fmt.Fprintf(file, "\t%s %s\n", fieldParts[0], fieldParts[1])
		}
	}

	file.WriteString("}\n\n")

	alias := strings.ToLower(string(className[0]))
	fmt.Fprintf(file, "func (%s *%s) Accept(visitor AstVisitor) interface {} {\n", alias, className)
	fmt.Fprintf(file, "\t return visitor.Visit%sExpr(%s)\n", className, alias)

	file.WriteString("}\n\n")
}

func defineVisitor(file *os.File, types []string) {
	file.WriteString("type AstVisitor interface {\n")

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])

		fmt.Fprintf(file, "\tVisit%sExpr (expt *%s) interface{}\n", className, className)
	}

	file.WriteString("}\n\n")
}
