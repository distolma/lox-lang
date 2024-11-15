package ast

import (
	"fmt"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(statements []Stmt) string {
	var result string
	for _, stmt := range statements {
		result += stmt.Accept(p).(string) + "\n"
	}
	return result
}

func (p *AstPrinter) PrintExpression(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *Variable) interface{} {
	return expr.Name.Lexeme
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}

	literal := expr.Value

	if v, ok := literal.(float64); ok {
		if v == float64(int(v)) {
			literal = fmt.Sprintf("%.1f", v)
		} else {
			literal = fmt.Sprintf("%g", v)
		}
	}

	return fmt.Sprintf("%v", literal)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitExpressionStmt(stmt *Expression) interface{} {
	return p.parenthesize("expr", stmt.Expression)
}

func (p *AstPrinter) VisitPrintStmt(stmt *Print) interface{} {
	return p.parenthesize("print", stmt.Expression)
}

func (p *AstPrinter) VisitVarStmt(stmt *Var) interface{} {
	if stmt.Initializer != nil {
		return p.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer)
	}
	return "(var " + stmt.Name.Lexeme + ")"
}

func (p *AstPrinter) execute(stmt Stmt) {
	stmt.Accept(p)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	result := "(" + name
	for _, expr := range exprs {
		result += " " + expr.Accept(p).(string)
	}
	result += ")"
	return result
}
