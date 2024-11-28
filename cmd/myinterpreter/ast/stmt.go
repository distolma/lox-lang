package ast

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
	VisitBlockStmt(expt *Block) interface{}
	VisitExpressionStmt(expt *Expression) interface{}
	VisitIfStmt(expt *If) interface{}
	VisitPrintStmt(expt *Print) interface{}
	VisitVarStmt(expt *Var) interface{}
	VisitWhileStmt(expt *While) interface{}
}

type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(b)
}

type Expression struct {
	Expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIfStmt(i)
}

type Print struct {
	Expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

type Var struct {
	Initializer Expr
	Name Token
}

func (v *Var) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(v)
}

type While struct {
	Condition Expr
	Body Stmt
}

func (w *While) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitWhileStmt(w)
}

