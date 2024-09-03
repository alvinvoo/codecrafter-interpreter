package scanner

import (
	"fmt"
)

type TokenType int

const (
	// Single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	VAR

	EOF
)

func (t TokenType) String() string {
	return [...]string{
		"LEFT_PAREN", "RIGHT_PAREN",
		"LEFT_BRACE", "RIGHT_BRACE",
		"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "SLASH", "STAR",
		"BANG", "BANG_EQUAL", "EQUAL", "EQUAL_EQUAL", "GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL",
		"IDENTIFIER", "STRING", "NUMBER",
		"VAR",
		"EOF",
	}[t]
}

type Scanner struct {
	source  []byte
	tokens  []string
	errors  []string
	start   int
	current int
	line    int
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []string{},
		errors:  []string{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) addToken(tokenType TokenType, lexeme string, literal string) {
	s.tokens = append(s.tokens, fmt.Sprintf("%s %s %s", tokenType.String(), lexeme, literal))
}

func (s *Scanner) addError(message string) {
	s.errors = append(s.errors, fmt.Sprintf("[line %d] Error: %s", s.line, message))
}

func (s *Scanner) nextMatch(expected byte) bool {
	if s.current == (len(s.source) - 1) {
		return false
	}

	if s.source[s.current+1] != expected {
		return false
	}

	return true
}

func (s *Scanner) advance() {
	s.current += 1
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) isLineBreak() bool {
	return s.source[s.current] == '\n'
}

func (s *Scanner) addEOF() {
	s.tokens = append(s.tokens, "EOF  null")
}

func (s *Scanner) addLine() {
	s.line += 1
}

func (s *Scanner) Tokenize() {
	for s.current < len(s.source) {
		t := s.source[s.current]

		switch t {
		case '(':
			s.addToken(LEFT_PAREN, string(t), "null")
		case ')':
			s.addToken(RIGHT_PAREN, string(t), "null")
		case '{':
			s.addToken(LEFT_BRACE, string(t), "null")
		case '}':
			s.addToken(RIGHT_BRACE, string(t), "null")
		case ',':
			s.addToken(COMMA, string(t), "null")
		case '.':
			s.addToken(DOT, string(t), "null")
		case '-':
			s.addToken(MINUS, string(t), "null")
		case '+':
			s.addToken(PLUS, string(t), "null")
		case ';':
			s.addToken(SEMICOLON, string(t), "null")
		case '*':
			s.addToken(STAR, string(t), "null")
		case '=':
			if s.nextMatch('=') {
				s.addToken(EQUAL_EQUAL, "==", "null")
				s.advance()
			} else {
				s.addToken(EQUAL, "==", "null")
			}
		case '!':
			if s.nextMatch('=') {
				s.addToken(BANG_EQUAL, "!=", "null")
				s.advance()
			} else {
				s.addToken(BANG, string(t), "null")
			}
		case '<':
			if s.nextMatch('=') {
				s.addToken(LESS_EQUAL, "<=", "null")
				s.advance()
			} else {
				s.addToken(LESS, string(t), "null")
			}
		case '>':
			if s.nextMatch('=') {
				s.addToken(GREATER_EQUAL, ">=", "null")
				s.advance()
			} else {
				s.addToken(GREATER, string(t), "null")
			}
		case '/':
			if s.nextMatch('/') {
				// it's a comment, ignore the rest of the line
				for !s.isAtEnd() && !s.isLineBreak() {
					s.advance()
				}
			} else {
				s.addToken(SLASH, string(t), "null")
			}
		case ' ': //ignore whitespace
		case '\t': //ignore tab
		case '\r': //ignore carriage returns
		case '\n':
			s.addLine()
		default:
			s.addError(fmt.Sprintf("Unexpected character: %s", string(t)))
		}

		s.advance()
	}

	s.addEOF()
}

func (s *Scanner) GetTokens() []string {
	return s.tokens
}

func (s *Scanner) GetErrors() []string {
	return s.errors
}
