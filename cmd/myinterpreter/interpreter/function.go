package interpreter

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	"github.com/distolma/golox/cmd/myinterpreter/environment"
)

type Callable interface {
	arity() int
	call(interpreter *Interpreter, arguments []interface{}) interface{}
}

type Function struct {
	declaraton ast.Function
	closure    *environment.Environment
}

func NewFunction(declaraton ast.Function, env *environment.Environment) *Function {
	return &Function{declaraton: declaraton, closure: env}
}

func (f *Function) arity() int {
	return len(f.declaraton.Params)
}

func (f *Function) call(interpreter *Interpreter, arguments []interface{}) (returnValue interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(Return); ok {
				returnValue = v.Value
				return
			}
			panic(err)
		}
	}()

	callEnv := environment.NewEnvironment(f.closure)
	for i, param := range f.declaraton.Params {
		callEnv.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaraton.Body, callEnv)
	return nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaraton.Name.Lexeme)
}
