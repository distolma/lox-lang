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
}

func NewFunction(declaraton ast.Function) *Function {
	return &Function{declaraton}
}

func (f *Function) arity() int {
	return len(f.declaraton.Params)
}

func (f *Function) call(interpreter *Interpreter, arguments []interface{}) (returnValue interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(Return); ok {

				fmt.Println(v.Value)
				returnValue = v.Value
				return
			}
			panic(err)
		}
	}()

	callEnv := environment.NewEnvironment(interpreter.globals)
	for i, param := range f.declaraton.Params {
		callEnv.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaraton.Body, callEnv)
	return nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaraton.Name.Lexeme)
}
