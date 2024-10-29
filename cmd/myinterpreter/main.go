package main

import (
	"bufio"
	"fmt"
	"os"

	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
	"github.com/distolma/golox/cmd/myinterpreter/scanner"
)

var log = &logerror.LogError{}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage

	filename := os.Args[2]
	runFile(filename)
}

func runPrompt() {
	inputScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !inputScanner.Scan() {
			break
		}

		line := inputScanner.Text()
		run(line)
		log.HadError = false
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	run(string(file))

	if log.HadError {
		os.Exit(65)
	}
}

func run(source string) {
	scan := scanner.NewScanner(source, log)
	tokens := scan.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
