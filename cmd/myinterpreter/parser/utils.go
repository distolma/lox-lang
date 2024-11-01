package parser

import (
	"github.com/distolma/golox/cmd/myinterpreter/ast"
)

func (p *Parser) match(tokens ...ast.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.EOF
}

func (p *Parser) peek() ast.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(t ast.TokenType, message string) ast.Token {
	if p.check(t) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return ast.Token{}
}
