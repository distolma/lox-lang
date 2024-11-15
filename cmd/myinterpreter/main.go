package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	"github.com/distolma/golox/cmd/myinterpreter/interpreter"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
	"github.com/distolma/golox/cmd/myinterpreter/parser"
	"github.com/distolma/golox/cmd/myinterpreter/scanner"
)

const (
	ExitError            = 1
	ExitCodeUsage        = 64
	ExitCodeSyntaxError  = 65
	ExitCodeRuntimeError = 70
)

type Lox struct {
	log         *logerror.LogError
	interpreter *interpreter.Interpreter
}

func NewLox() *Lox {
	log := &logerror.LogError{}
	interpreter := interpreter.NewInterpreter(log)

	return &Lox{
		log:         log,
		interpreter: interpreter,
	}
}

func main() {
	lox := NewLox()

	if len(os.Args) < 3 {
		lox.runPrompt()
		return
	}

	command := os.Args[1]
	filename := os.Args[2]

	validCommands := []string{"tokenize", "parse", "evaluate", "run"}
	if slices.Contains(validCommands, command) {
		switch command {
		case "tokenize":
			lox.tokenize(filename)
		case "parse":
			lox.parse(filename)
		case "evaluate":
			lox.evaluate(filename)
		case "run":
			lox.runFile(filename)
		default:
			lox.runFile(filename)
		}
		return
	}
}

func (l *Lox) runPrompt() {
	inputScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !inputScanner.Scan() {
			break
		}
		line := inputScanner.Text()
		l.run(line)
		l.log.HadError = false
	}
}

func (l *Lox) runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(ExitError)
	}
	l.run(string(file))

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	if l.log.HadRuntimeError {
		os.Exit(ExitCodeRuntimeError)
	}
}

func (l *Lox) run(source string) {
	scanner := scanner.NewScanner(source, l.log)
	tokens := scanner.ScanTokens()

	parser := parser.NewParser(tokens, l.log)
	statements := parser.Parse()

	if l.log.HadError {
		return
	}

	l.interpreter.Interpret(statements)
}

func (l *Lox) tokenize(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(ExitError)
	}
	source := string(file)

	scan := scanner.NewScanner(source, l.log)
	tokens := scan.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.String())
	}

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}
}

func (l *Lox) parse(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(ExitError)
	}
	source := string(file)

	scan := scanner.NewScanner(source, l.log)
	tokens := scan.ScanTokens()

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	parser := parser.NewParser(tokens, l.log)
	expression := parser.ParseExpression()

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	printer := ast.AstPrinter{}
	result := printer.PrintExpression(expression)
	fmt.Println(result)
}

func (l *Lox) evaluate(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(ExitError)
	}
	source := string(file)

	scan := scanner.NewScanner(source, l.log)
	tokens := scan.ScanTokens()

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	parser := parser.NewParser(tokens, l.log)
	expression := parser.ParseExpression()

	value := l.interpreter.InterpretExpression(expression)
	if l.log.HadRuntimeError {
		os.Exit(ExitCodeRuntimeError)
	}
	fmt.Println(value)
}
