package parser

import (
	"fmt"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
)

type Parser struct {
	log     *logerror.LogError
	tokens  []ast.Token
	current int
}

func NewParser(tokens []ast.Token, log *logerror.LogError) *Parser {
	return &Parser{tokens: tokens, current: 0, log: log}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(ParseError); ok {
				p.synchronize()
			} else {
				panic(err)
			}
		}
	}()

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) ParseExpression() ast.Expr {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(ParseError); ok {
				p.synchronize()
			} else {
				panic(err)
			}
		}
	}()

	return p.expression()
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() ast.Stmt {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(ParseError); ok {
				p.synchronize()
			} else {
				panic(err)
			}
		}
	}()

	if p.match(ast.TFun) {
		return p.function("function")
	}

	if p.match(ast.TVar) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) statement() ast.Stmt {
	if p.match(ast.TFor) {
		return p.forStatement()
	}
	if p.match(ast.TIf) {
		return p.ifStatement()
	}
	if p.match(ast.TPrint) {
		return p.printStatement()
	}
	if p.match(ast.TReturn) {
		return p.returnStatement()
	}
	if p.match(ast.TWhile) {
		return p.whileStatement()
	}
	if p.match(ast.TLeftBrace) {
		return &ast.Block{Statements: p.block()}
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(ast.TLeftParen, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(ast.TSemicolon) {
		initializer = nil
	} else if p.match(ast.TVar) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(ast.TSemicolon) {
		condition = p.expression()
	}
	p.consume(ast.TSemicolon, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(ast.TRightParen) {
		increment = p.expression()
	}
	p.consume(ast.TRightParen, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &ast.Block{Statements: []ast.Stmt{body, &ast.Expression{Expression: increment}}}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.While{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.Block{Statements: []ast.Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(ast.TLeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(ast.TRightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(ast.TElse) {
		elseBranch = p.statement()
	}

	return &ast.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(ast.TSemicolon, "Expect ';' after value.")
	return &ast.Print{Expression: expr}
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr
	if !p.check(ast.TSemicolon) {
		value = p.expression()
	}
	p.consume(ast.TSemicolon, "Expect ';' after return value.")
	return &ast.Return{Keyword: keyword, Value: value}
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(ast.TIdentifier, "Expect variable name.")

	var initializer ast.Expr
	if p.match(ast.TEqual) {
		initializer = p.expression()
	}

	p.consume(ast.TSemicolon, "Expect ';' after variable declaration.")
	return &ast.Var{Name: name, Initializer: initializer}
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(ast.TLeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(ast.TRightParen, "Expect ')' after condition.")
	body := p.statement()

	return &ast.While{Condition: condition, Body: body}
}

func (p *Parser) expressionStatement() ast.Stmt {
	value := p.expression()
	p.consume(ast.TSemicolon, "Expect ';' after expression.")
	return &ast.Expression{Expression: value}
}

func (p *Parser) function(kind string) *ast.Function {
	name := p.consume(ast.TIdentifier, fmt.Sprintf("Expect %s name.", kind))

	p.consume(ast.TLeftParen, fmt.Sprintf("Expact '(' after %s name.", kind))

	var parameters []ast.Token
	if !p.check(ast.TRightParen) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}

			paramToken := p.consume(ast.TIdentifier, "Expect parameter name.")

			parameters = append(parameters, paramToken)

			if !p.match(ast.TComma) {
				break
			}
		}
	}

	p.consume(ast.TRightParen, "Expect ')' after parameters.")

	p.consume(ast.TLeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()

	return &ast.Function{Name: name, Params: parameters, Body: body}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt

	for !p.check(ast.TRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(ast.TRightBrace, "Expect '}' after block.")

	return statements
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(ast.TEqual) {
		equals := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(*ast.Variable); ok {
			return &ast.Assign{Name: varExpr.Name, Value: value}
		}

		p.error(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(ast.TOr) {
		operator := p.previous()
		right := p.and()
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(ast.TAnd) {
		operator := p.previous()
		right := p.equality()
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(ast.TBangEqual, ast.TEqualEqual) {
		operator := p.previous()
		right := p.comparison()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(ast.TGreater, ast.TGreaterEqual, ast.TLess, ast.TLessEqual) {
		operator := p.previous()
		right := p.term()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(ast.TPlus, ast.TMinus) {
		operator := p.previous()
		right := p.factor()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(ast.TSlash, ast.TStar) {
		operator := p.previous()
		right := p.unary()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(ast.TBang, ast.TMinus) {
		operator := p.previous()
		right := p.unary()

		return &ast.Unary{Operator: operator, Right: right}

	}

	return p.call()
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr

	if !p.check(ast.TRightParen) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())

			if !p.match(ast.TComma) {
				break
			}
		}
	}

	paren := p.consume(ast.TRightParen, "Expect ')' after arguments.")

	return &ast.Call{Callee: callee, Arguments: arguments, Paren: paren}
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(ast.TLeftParen) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) primary() ast.Expr {
	if p.match(ast.TFalse) {
		return &ast.Literal{Value: false}
	} else if p.match(ast.TTrue) {
		return &ast.Literal{Value: true}
	} else if p.match(ast.TNil) {
		return &ast.Literal{Value: nil}
	} else if p.match(ast.TNumber, ast.TString) {
		return &ast.Literal{Value: p.previous().Literal}
	} else if p.match(ast.TIdentifier) {
		return &ast.Variable{Name: p.previous()}
	} else if p.match(ast.TLeftParen) {
		expr := p.expression()
		p.consume(ast.TRightParen, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == ast.TSemicolon {
			return
		}

		switch p.peek().Type {
		case ast.TClass:
		case ast.TFun:
		case ast.TVar:
		case ast.TFor:
		case ast.TIf:
		case ast.TWhile:
		case ast.TPrint:
		case ast.TReturn:
			return
		}

		p.advance()
	}
}
