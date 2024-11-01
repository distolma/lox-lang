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

	for p.match(ast.BangEqual, ast.EqualEqual) {
		operator := p.previous()
		right := p.comparison()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(ast.Greater, ast.GreaterEqual, ast.Less, ast.LessEqual) {
		operator := p.previous()
		right := p.term()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(ast.Plus, ast.Minus) {
		operator := p.previous()
		right := p.factor()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(ast.Slash, ast.Star) {
		operator := p.previous()
		right := p.unary()

		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(ast.Bang, ast.Minus) {
		operator := p.previous()
		right := p.unary()

		return &ast.Unary{Operator: operator, Right: right}

	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(ast.False) {
		return &ast.Literal{Value: false}
	} else if p.match(ast.True) {
		return &ast.Literal{Value: true}
	} else if p.match(ast.Nil) {
		return &ast.Literal{Value: nil}
	} else if p.match(ast.Number, ast.String) {
		return &ast.Literal{Value: p.previous().Literal}
	} else if p.match(ast.LeftParen) {
		expr := p.expression()
		p.consume(ast.RightParen, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == ast.Semicolon {
			return
		}

		switch p.peek().Type {
		case ast.Class:
		case ast.Fun:
		case ast.Var:
		case ast.For:
		case ast.If:
		case ast.While:
		case ast.Print:
		case ast.Return:
			return
		}

		p.advance()
	}
}
