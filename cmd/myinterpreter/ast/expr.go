package ast

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(expt *Binary) interface{}
	VisitGroupingExpr(expt *Grouping) interface{}
	VisitLiteralExpr(expt *Literal) interface{}
	VisitUnaryExpr(expt *Unary) interface{}
}

type Binary struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Right    Expr
	Operator Token
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}
