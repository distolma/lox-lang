package interpreter

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	"github.com/distolma/golox/cmd/myinterpreter/environment"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
)

type Interpreter struct {
	log         *logerror.LogError
	environment *environment.Environment
	globals     *environment.Environment
}

func NewInterpreter(log *logerror.LogError) *Interpreter {
	globals := environment.NewEnvironment(nil)
	globals.Define("clock", Clock{})

	return &Interpreter{
		log:         log,
		environment: globals,
		globals:     globals,
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		if err := recover(); err != nil {
			if runtimeError, ok := err.(RuntimeError); ok {
				i.log.RuntimeError(runtimeError.Token, runtimeError.Message)
			} else {
				panic(err)
			}
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) InterpretExpression(expr ast.Expr) string {
	defer func() {
		if err := recover(); err != nil {
			if runtimeError, ok := err.(RuntimeError); ok {
				i.log.RuntimeError(runtimeError.Token, runtimeError.Message)
			} else {
				panic(err)
			}
		}
	}()

	return i.stringify(i.evaluate(expr))
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == ast.TOr {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.TBang:
		return !i.isTruthy(right)
	case ast.TMinus:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}

	// unreachable
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) interface{} {
	value, err := i.environment.Get(expr.Name.Lexeme)
	if err != nil {
		panic(NewRuntimeError(expr.Name, err.Error()))
	}
	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.TMinus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case ast.TSlash:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case ast.TStar:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			panic(NewRuntimeError(expr.Operator, "Division by zero."))
		}
		return left.(float64) * right.(float64)
	case ast.TPlus:
		leftFloat, leftOk := left.(float64)
		rightFloat, rightOk := right.(float64)

		if leftOk && rightOk {
			return leftFloat + rightFloat
		}

		leftString, leftOk := left.(string)
		rightString, rightOk := right.(string)
		if leftOk && rightOk {
			return leftString + rightString
		}

		panic(NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings."))
	case ast.TGreater:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case ast.TGreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case ast.TLess:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case ast.TLessEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case ast.TBangEqual:
		return left != right
	case ast.TEqualEqual:
		return left == right
	}

	// unreachable
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.Call) interface{} {
	callee := i.evaluate(expr.Callee)

	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := (callee).(Callable)
	if !ok {
		panic(NewRuntimeError(expr.Paren, "Can only call functions and classes."))
	}

	if len(arguments) != function.arity() {
		panic(NewRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))))
	}

	return function.call(i, arguments)
}

func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *environment.Environment) {
	previous := i.environment

	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) interface{} {
	i.executeBlock(stmt.Statements, environment.NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) interface{} {
	function := NewFunction(*stmt, i.environment)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) interface{} {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)
	err := i.environment.Assign(expr.Name.Lexeme, value)
	if err != nil {
		panic(NewRuntimeError(expr.Name, err.Error()))
	}
	return value
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if v, ok := object.(bool); ok {
		return v
	}

	return true
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	return fmt.Sprint(object)
}

func (i *Interpreter) checkNumberOperand(operator ast.Token, operand interface{}) {
	if _, ok := operand.(float64); ok {
		return
	}

	panic(NewRuntimeError(operator, "Operand must be a number."))
}

func (i *Interpreter) checkNumberOperands(operator ast.Token, left interface{}, right interface{}) {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return
	}

	panic(NewRuntimeError(operator, "Operands must be numbers."))
}
