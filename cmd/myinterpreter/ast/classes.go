package ast

type AstVisitor interface {
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

func (b *Binary) Accept(visitor AstVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor AstVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor AstVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Right    Expr
	Operator Token
}

func (u *Unary) Accept(visitor AstVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}
