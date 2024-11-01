package ast

type Expr interface {
	Accept(visitor AstVisitor) interface{}
}
