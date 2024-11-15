package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	"github.com/distolma/golox/cmd/myinterpreter/interpreter"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
	"github.com/distolma/golox/cmd/myinterpreter/parser"
	"github.com/distolma/golox/cmd/myinterpreter/scanner"
)

var log = &logerror.LogError{}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" && command != "evaluate" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage

	filename := os.Args[2]
	runFile(filename, command)
}

func runPrompt(command string) {
	inputScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !inputScanner.Scan() {
			break
		}

		line := inputScanner.Text()
		run(line, command)
		log.HadError = false
	}
}

func runFile(path string, command string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	run(string(file), command)
}

func run(source string, command string) {
	scan := scanner.NewScanner(source, log)
	tokens := scan.ScanTokens()

	if command == "tokenize" {
		for _, token := range tokens {
			fmt.Println(token.String())
		}

		if log.HadError {
			os.Exit(65)
		}
		return
	}

	if log.HadError {
		os.Exit(65)
	}

	parser := parser.NewParser(tokens, log)
	expression := parser.Parse()

	if log.HadError {
		os.Exit(65)
	}

	if command == "parse" {
		printer := ast.AstPrinter{}
		result := printer.Print(expression)
		fmt.Println(result)
		return
	}

	interpreter := interpreter.NewInterpreter()

	if command == "evaluate" {
		value := interpreter.Interpret(expression)

		fmt.Println(value)
	}
}
