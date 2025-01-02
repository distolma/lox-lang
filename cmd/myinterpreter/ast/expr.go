package ast

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitAssignExpr(expt *Assign) interface{}
	VisitBinaryExpr(expt *Binary) interface{}
	VisitCallExpr(expt *Call) interface{}
	VisitGroupingExpr(expt *Grouping) interface{}
	VisitLiteralExpr(expt *Literal) interface{}
	VisitLogicalExpr(expt *Logical) interface{}
	VisitUnaryExpr(expt *Unary) interface{}
	VisitVariableExpr(expt *Variable) interface{}
}

type Assign struct {
	Value Expr
	Name  Token
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

type Binary struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (c *Call) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(c)
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

type Logical struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(l)
}

type Unary struct {
	Right    Expr
	Operator Token
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

type Variable struct {
	Name Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}
