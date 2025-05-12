package resolver

import (
	"github.com/distolma/golox/cmd/myinterpreter/ast"
	"github.com/distolma/golox/cmd/myinterpreter/interpreter"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
)

type Scope map[string]bool

func (s *Scope) Declare(name string) {
	(*s)[name] = false
}

func (s *Scope) Define(name string) {
	(*s)[name] = true
}

func (s *Scope) Has(name string) (declared, defined bool) {
	if value, ok := (*s)[name]; ok {
		return true, value
	}
	return false, false
}

type Stack []Scope

func (s *Stack) Peek() *Scope {
	return &(*s)[s.Size()-1]
}

func (s *Stack) Push(value Scope) {
	*s = append(*s, value)
}

func (s *Stack) Pop() {
	*s = (*s)[:s.Size()-1]
}

func (s *Stack) Size() int {
	return len(*s)
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

const (
	FunctionTypeNone = iota
	FunctionTypeFunction
)

type Resolver struct {
	log             *logerror.LogError
	interpreter     *interpreter.Interpreter
	scopes          Stack
	currentFunction int
}

func NewResolver(interpreter *interpreter.Interpreter, log *logerror.LogError) *Resolver {
	return &Resolver{interpreter: interpreter, log: log, currentFunction: FunctionTypeNone}
}

func (r *Resolver) ResolveStmts(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(statement ast.Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(function *ast.Function, functionType int) {
	enclosingFunction := r.currentFunction
	r.beginScope()
	r.currentFunction = functionType
	defer func() {
		r.endScope()
		r.currentFunction = enclosingFunction
	}()

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.ResolveStmts(function.Body)
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(Scope))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name ast.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	if _, defined := r.scopes.Peek().Has(name.Lexeme); defined {
		r.log.TokenError(name, "Already a variable with this name in this scope.")
	}
	r.scopes.Peek().Declare(name.Lexeme)
}

func (r *Resolver) define(name ast.Token) {
	if r.scopes.IsEmpty() {
		return
	}
	r.scopes.Peek().Define(name.Lexeme)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name ast.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, defined := r.scopes[i].Has(name.Lexeme); defined {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) VisitBlockStmt(stmt *ast.Block) interface{} {
	r.beginScope()
	r.ResolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.Expression) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.Function) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FunctionTypeFunction)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.If) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.Print) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.Return) interface{} {
	if r.currentFunction == FunctionTypeNone {
		r.log.TokenError(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.Var) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.While) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)

	return nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.Call) interface{} {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.Grouping) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.Literal) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.Logical) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.Variable) interface{} {
	if !r.scopes.IsEmpty() {
		if declared, defined := r.scopes.Peek().Has(expr.Name.Lexeme); declared && !defined {
			r.log.TokenError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil
}
