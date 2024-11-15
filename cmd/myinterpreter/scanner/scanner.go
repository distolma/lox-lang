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
		s.addToken(ast.TLeftParen)
	case ')':
		s.addToken(ast.TRightParen)
	case '{':
		s.addToken(ast.TLeftBrace)
	case '}':
		s.addToken(ast.TRightBrace)
	case ',':
		s.addToken(ast.TComma)
	case '.':
		s.addToken(ast.TDot)
	case '-':
		s.addToken(ast.TMinus)
	case '+':
		s.addToken(ast.TPlus)
	case ';':
		s.addToken(ast.TSemicolon)
	case '*':
		s.addToken(ast.TStar)
	case '!':
		if s.match('=') {
			s.addToken(ast.TBangEqual)
		} else {
			s.addToken(ast.TBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.TEqualEqual)
		} else {
			s.addToken(ast.TEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.TLessEqual)
		} else {
			s.addToken(ast.TLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.TGreaterEqual)
		} else {
			s.addToken(ast.TGreater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.TSlash)
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
	s.addTokenWithLiteral(ast.TString, value)
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
	s.addTokenWithLiteral(ast.TNumber, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, found := keywords[text]
	if !found {
		tokenType = ast.TIdentifier
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
	"and":    ast.TAnd,
	"class":  ast.TClass,
	"else":   ast.TElse,
	"false":  ast.TFalse,
	"for":    ast.TFor,
	"fun":    ast.TFun,
	"if":     ast.TIf,
	"nil":    ast.TNil,
	"or":     ast.TOr,
	"print":  ast.TPrint,
	"return": ast.TReturn,
	"super":  ast.TSuper,
	"this":   ast.TThis,
	"true":   ast.TTrue,
	"var":    ast.TVar,
	"while":  ast.TWhile,
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
