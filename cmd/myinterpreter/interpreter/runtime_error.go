package interpreter

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
)

type RuntimeError struct {
	message string
	token   ast.Token
}

func NewRuntimeError(token ast.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", re.message, re.token.Line)
}
