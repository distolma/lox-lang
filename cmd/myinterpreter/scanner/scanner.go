package scanner

import (
	"fmt"
	"strconv"

	"github.com/distolma/golox/cmd/myinterpreter/ast"
	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
)

type Scanner struct {
	source  string
	log     *logerror.LogError
	tokens  []ast.Token
	current int
	line    int
	start   int
}

func NewScanner(source string, log *logerror.LogError) *Scanner {
	return &Scanner{source: source, line: 1, log: log}
}

func (s *Scanner) ScanTokens() []ast.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, ast.Token{Type: ast.EOF, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(ast.LeftParen)
	case ')':
		s.addToken(ast.RightParen)
	case '{':
		s.addToken(ast.LeftBrace)
	case '}':
		s.addToken(ast.RightBrace)
	case ',':
		s.addToken(ast.Comma)
	case '.':
		s.addToken(ast.Dot)
	case '-':
		s.addToken(ast.Minus)
	case '+':
		s.addToken(ast.Plus)
	case ';':
		s.addToken(ast.Semicolon)
	case '*':
		s.addToken(ast.Star)
	case '!':
		if s.match('=') {
			s.addToken(ast.BangEqual)
		} else {
			s.addToken(ast.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EqualEqual)
		} else {
			s.addToken(ast.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LessEqual)
		} else {
			s.addToken(ast.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GreaterEqual)
		} else {
			s.addToken(ast.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.Slash)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(char) {
			s.number()
		} else if s.isAlpha(char) {
			s.identifier()
		} else {
			s.log.Error(s.line, fmt.Sprintf("Unexpected character: %s", string(char)))
		}
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.log.Error(s.line, "Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(ast.String, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenWithLiteral(ast.Number, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = ast.Identifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) addToken(tokenType ast.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType ast.TokenType, literal interface{}) {
	token := ast.Token{
		Lexeme:  s.source[s.start:s.current],
		Type:    tokenType,
		Literal: literal,
		Line:    s.line,
	}
	s.tokens = append(s.tokens, token)
}

var keywords = map[string]ast.TokenType{
	"and":    ast.And,
	"class":  ast.Class,
	"else":   ast.Else,
	"false":  ast.False,
	"for":    ast.For,
	"fun":    ast.Fun,
	"if":     ast.If,
	"nil":    ast.Nil,
	"or":     ast.Or,
	"print":  ast.Print,
	"return": ast.Return,
	"super":  ast.Super,
	"this":   ast.This,
	"true":   ast.True,
	"var":    ast.Var,
	"while":  ast.While,
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	char := rune(s.source[s.current])
	s.current++
	return char
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) isAlpha(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func (s *Scanner) isAlphaNumeric(char rune) bool {
	return s.isDigit(char) || s.isAlpha(char)
}
