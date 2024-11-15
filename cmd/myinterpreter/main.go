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
	interpreter *interpreter.Interpreter
	log         *logerror.LogError
}

func NewLox() *Lox {
	log := &logerror.LogError{}

	return &Lox{
		log:         log,
		interpreter: interpreter.NewInterpreter(log),
	}
}

func main() {
	lox := NewLox()

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(ExitError)
	}

	command := os.Args[1]
	filename := os.Args[2]

	validCommands := []string{"tokenize", "parse", "evaluate", "run"}
	if slices.Contains(validCommands, command) {
		lox.runFile(filename, command)
		return
	}

	lox.runPrompt()
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

func (l *Lox) runFile(path string, command string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(ExitError)
	}
	l.runCommand(string(file), command)

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	if l.log.HadRuntimeError {
		os.Exit(ExitCodeRuntimeError)
	}
}

func (l *Lox) runCommand(source string, command string) {
	scan := scanner.NewScanner(source, l.log)
	tokens := scan.ScanTokens()

	if l.log.HadError {
		os.Exit(ExitCodeSyntaxError)
	}

	if command == "tokenize" {
		for _, token := range tokens {
			fmt.Println(token.String())
		}
		return
	}

	parser := parser.NewParser(tokens, l.log)
	expression := parser.Parse()

	if l.log.HadError || expression == nil {
		os.Exit(ExitCodeSyntaxError)
	}

	if command == "parse" {
		printer := ast.AstPrinter{}
		result := printer.Print(expression)
		fmt.Println(result)
		return
	}

	if command == "evaluate" {
		value := l.interpreter.Interpret(expression)
		if l.log.HadRuntimeError {
			os.Exit(ExitCodeRuntimeError)
		}
		fmt.Println(value)
	}
}

func (l *Lox) run(source string) {
	scanner := scanner.NewScanner(source, l.log)
	tokens := scanner.ScanTokens()

	parser := parser.NewParser(tokens, l.log)
	expression := parser.Parse()

	if l.log.HadError {
		return
	}

	l.interpreter.Interpret(expression)
}
