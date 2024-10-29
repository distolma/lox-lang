package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/distolma/golox/cmd/myinterpreter/scanner"
)

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
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		run(string(fileContents))
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
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
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	run(string(file))
}

func run(source string) {
	scan := scanner.NewScanner(source)
	tokens := scan.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
