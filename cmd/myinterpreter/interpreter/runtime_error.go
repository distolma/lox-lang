package interpreter

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
)

type RuntimeError struct {
	Message string
	Token   ast.Token
}

func NewRuntimeError(token ast.Token, message string) RuntimeError {
	return RuntimeError{Token: token, Message: message}
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", re.Message, re.Token.Line)
}
