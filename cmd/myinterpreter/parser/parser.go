package parser

import (
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

func (p *Parser) Parse() ast.Expr {
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
	return p.equality()
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

	return p.primary()
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
