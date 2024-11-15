package interpreter

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
)

type Interpreter struct {
	log *logerror.LogError
}

func NewInterpreter(log *logerror.LogError) *Interpreter {
	return &Interpreter{log: log}
}

func (i *Interpreter) Interpret(expr ast.Expr) string {
	defer func() {
		if err := recover(); err != nil {
			if runtimeError, ok := err.(RuntimeError); ok {
				i.log.RuntimeError(runtimeError.Token, runtimeError.Message)
			} else {
				// fmt.Println(err)
				// panic(err)
			}
		}
	}()

	return i.stringify(i.evaluate(expr))
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.Bang:
		return !i.isTruthy(right)
	case ast.Minus:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}

	// unreachable
	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case ast.Minus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case ast.Slash:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case ast.Star:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			panic(NewRuntimeError(expr.Operator, "Division by zero."))
		}
		return left.(float64) * right.(float64)
	case ast.Plus:
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
	case ast.Greater:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case ast.GreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case ast.Less:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case ast.LessEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case ast.BangEqual:
		return left != right
	case ast.EqualEqual:
		return left == right
	}

	// unreachable
	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(i)
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
