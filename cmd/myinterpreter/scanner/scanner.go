package scanner

import (
	"fmt"
	"strconv"

	logerror "github.com/distolma/golox/cmd/myinterpreter/log_error"
	"github.com/distolma/golox/cmd/myinterpreter/token"
)

type Scanner struct {
	source  string
	log     *logerror.LogError
	tokens  []token.Token
	current int
	line    int
	start   int
}

func NewScanner(source string, log *logerror.LogError) *Scanner {
	return &Scanner{source: source, line: 1, log: log}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.Token{Type: token.EOF, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(token.LeftParen)
	case ')':
		s.addToken(token.RightParen)
	case '{':
		s.addToken(token.LeftBrace)
	case '}':
		s.addToken(token.RightBrace)
	case ',':
		s.addToken(token.Comma)
	case '.':
		s.addToken(token.Dot)
	case '-':
		s.addToken(token.Minus)
	case '+':
		s.addToken(token.Plus)
	case ';':
		s.addToken(token.Semicolon)
	case '*':
		s.addToken(token.Star)
	case '!':
		if s.match('=') {
			s.addToken(token.BangEqual)
		} else {
			s.addToken(token.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EqualEqual)
		} else {
			s.addToken(token.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LessEqual)
		} else {
			s.addToken(token.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GreaterEqual)
		} else {
			s.addToken(token.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.Slash)
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
	s.addTokenWithLiteral(token.String, value)
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
	s.addTokenWithLiteral(token.Number, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = token.Identifier
	}

	s.addToken(tokenType)
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
	token := token.Token{
		Lexeme:  s.source[s.start:s.current],
		Type:    tokenType,
		Literal: literal,
		Line:    s.line,
	}
	s.tokens = append(s.tokens, token)
}

var keywords = map[string]token.TokenType{
	"and":    token.And,
	"class":  token.Class,
	"else":   token.Else,
	"false":  token.False,
	"for":    token.For,
	"fun":    token.Fun,
	"if":     token.If,
	"nil":    token.Nil,
	"or":     token.Or,
	"print":  token.Print,
	"return": token.Return,
	"super":  token.Super,
	"this":   token.This,
	"true":   token.True,
	"var":    token.Var,
	"while":  token.While,
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
