package ast

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpressionStmt(expt *Expression) interface{}
	VisitPrintStmt(expt *Print) interface{}
}

type Expression struct {
	Expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

type Print struct {
	Expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}
