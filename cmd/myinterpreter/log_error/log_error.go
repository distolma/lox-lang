package logerror

import (
	"fmt"
	"os"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
)

type LogError struct {
	HadError        bool
	HadRuntimeError bool
}

func (l *LogError) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.HadError = true
}

func (l *LogError) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *LogError) TokenError(token ast.Token, message string) {
	var where string
	if token.Type == ast.EOF {
		where = " at end"
	} else {
		where = fmt.Sprintf(" at '%s'", token.Lexeme)
	}

	l.report(token.Line, where, message)
}

func (l *LogError) RuntimeError(token ast.Token, message string) {
	fmt.Fprintf(os.Stderr, "%s \n[line: %d]", message, token.Line)
	l.HadRuntimeError = true
}
