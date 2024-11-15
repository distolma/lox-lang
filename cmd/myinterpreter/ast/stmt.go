package ast

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitExpressionStmt(expt *Expression) interface{}
	VisitPrintStmt(expt *Print) interface{}
	VisitVarStmt(expt *Var) interface{}
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

type Var struct {
	Initializer Expr
	Name        Token
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}
