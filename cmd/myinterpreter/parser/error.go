package parser

import "github.com/distolma/golox/cmd/myinterpreter/ast"

type ParseError struct {
	msg string
}

func (p *ParseError) Error() string {
	return p.msg
}

func (p *Parser) error(token ast.Token, message string) {
	p.log.TokenError(token, message)
	panic(ParseError{})
}
