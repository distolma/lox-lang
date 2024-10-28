package scanner

import (
	"fmt"
	"strconv"

	"github.com/distolma/golox/cmd/myinterpreter/tokens"
)

type Scanner struct {
	source  string
	tokens  []tokens.Token
	current int
	line    int
	start   int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, line: 1}
}

func (s *Scanner) ScanTokens() []tokens.Token {
	if !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tokens.Token{Type: tokens.EOF, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(tokens.LeftParen)
	case ')':
		s.addToken(tokens.RightParen)
	case '{':
		s.addToken(tokens.LeftBrace)
	case '}':
		s.addToken(tokens.RightBrace)
	case ',':
		s.addToken(tokens.Comma)
	case '.':
		s.addToken(tokens.Dot)
	case '-':
		s.addToken(tokens.Minus)
	case '+':
		s.addToken(tokens.Plus)
	case ';':
		s.addToken(tokens.Semicolon)
	case '*':
		s.addToken(tokens.Star)
	case '!':
		if s.match('=') {
			s.addToken(tokens.BangEqual)
		} else {
			s.addToken(tokens.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(tokens.EqualEqual)
		} else {
			s.addToken(tokens.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(tokens.LessEqual)
		} else {
			s.addToken(tokens.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(tokens.GreaterEqual)
		} else {
			s.addToken(tokens.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tokens.Slash)
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
		} else {
			fmt.Print("Unexpected character.")
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
		fmt.Print("Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(tokens.String, value)
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
	s.addTokenWithLiteral(tokens.Number, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = tokens.Identifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) addToken(tokenType tokens.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType tokens.TokenType, literal interface{}) {
	token := tokens.Token{
		Lexeme:  s.source[s.start:s.current],
		Type:    tokenType,
		Literal: literal,
		Line:    s.line,
	}
	s.tokens = append(s.tokens, token)
}

var keywords = map[string]tokens.TokenType{
	"and":    tokens.And,
	"class":  tokens.Class,
	"else":   tokens.Else,
	"false":  tokens.False,
	"for":    tokens.For,
	"fun":    tokens.Fun,
	"if":     tokens.If,
	"nil":    tokens.Nil,
	"or":     tokens.Or,
	"print":  tokens.Print,
	"return": tokens.Return,
	"super":  tokens.Super,
	"this":   tokens.This,
	"true":   tokens.True,
	"var":    tokens.Var,
	"while":  tokens.While,
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
