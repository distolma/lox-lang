package ast

import (
	"fmt"
)

type TokenType string

const (
	// Single-character tokens
	TLeftParen  TokenType = "LEFT_PAREN"
	TRightParen TokenType = "RIGHT_PAREN"
	TLeftBrace  TokenType = "LEFT_BRACE"
	TRightBrace TokenType = "RIGHT_BRACE"
	TComma      TokenType = "COMMA"
	TDot        TokenType = "DOT"
	TMinus      TokenType = "MINUS"
	TPlus       TokenType = "PLUS"
	TSemicolon  TokenType = "SEMICOLON"
	TSlash      TokenType = "SLASH"
	TStar       TokenType = "STAR"
	// One or two character tokens
	TBang         TokenType = "BANG"
	TBangEqual    TokenType = "BANG_EQUAL"
	TEqual        TokenType = "EQUAL"
	TEqualEqual   TokenType = "EQUAL_EQUAL"
	TGreater      TokenType = "GREATER"
	TGreaterEqual TokenType = "GREATER_EQUAL"
	TLess         TokenType = "LESS"
	TLessEqual    TokenType = "LESS_EQUAL"
	// Literals
	TIdentifier TokenType = "IDENTIFIER"
	TString     TokenType = "STRING"
	TNumber     TokenType = "NUMBER"
	// Keywords
	TAnd    TokenType = "AND"
	TClass  TokenType = "CLASS"
	TElse   TokenType = "ELSE"
	TFalse  TokenType = "FALSE"
	TFun    TokenType = "FUN"
	TFor    TokenType = "FOR"
	TIf     TokenType = "IF"
	TNil    TokenType = "NIL"
	TOr     TokenType = "OR"
	TPrint  TokenType = "PRINT"
	TReturn TokenType = "RETURN"
	TSuper  TokenType = "SUPER"
	TThis   TokenType = "THIS"
	TTrue   TokenType = "TRUE"
	TVar    TokenType = "VAR"
	TWhile  TokenType = "WHILE"

	EOF TokenType = "EOF"
)

type Token struct {
	Literal interface{}
	Lexeme  string
	Type    TokenType
	Line    int
}

func (t *Token) String() string {
	literal := t.Literal

	if literal == nil {
		literal = "null"
	}

	if v, ok := t.Literal.(float64); ok {
		if v == float64(int(v)) {
			literal = fmt.Sprintf("%.1f", v)
		} else {
			literal = fmt.Sprintf("%g", v)
		}
	}
	return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, literal)
}
