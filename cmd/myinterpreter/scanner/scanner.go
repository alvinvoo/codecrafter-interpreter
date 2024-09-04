package scanner

import (
	"fmt"
	"math"
	"strconv"
	"unicode"
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
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR

	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (t TokenType) String() string {
	return [...]string{
		"LEFT_PAREN", "RIGHT_PAREN",
		"LEFT_BRACE", "RIGHT_BRACE",
		"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "SLASH", "STAR",
		"BANG", "BANG_EQUAL", "EQUAL", "EQUAL_EQUAL", "GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL",
		"IDENTIFIER", "STRING", "NUMBER",
		"AND", "CLASS", "ELSE", "FALSE", "FUN", "FOR", "IF", "NIL", "OR",
		"PRINT", "RETURN", "SUPER", "THIS", "TRUE", "VAR", "WHILE",
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
	// TODO: can do improvement of `start` and `current` to keep track of the position of the token
	s.tokens = append(s.tokens, fmt.Sprintf("%s %s %s", tokenType.String(), lexeme, literal))
}

func (s *Scanner) addError(message string) {
	s.errors = append(s.errors, fmt.Sprintf("[line %d] Error: %s", s.line, message))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)-1
}

func (s *Scanner) nextMatch(expected byte) bool {
	if s.isAtEnd() {
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

func (s *Scanner) addEOF() {
	s.tokens = append(s.tokens, "EOF  null")
}

func (s *Scanner) addLine() {
	s.line += 1
}

func (s *Scanner) addNumber() {
	numStr := string(s.source[s.current])
	for !s.isAtEnd() && (unicode.IsDigit(rune(s.source[s.current+1])) || s.nextMatch('.')) {
		s.advance()
		curDigit := s.source[s.current]

		numStr += string(curDigit)
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		s.addError(fmt.Sprintf("Error parsing number: %s", numStr))
	}

	var literal string
	// Separate the integer and fractional parts
	_, frac := math.Modf(num)

	// If the fractional part is not zero, it's a float
	if frac != 0 {
		literal = strconv.FormatFloat(num, 'f', -1, 64)
	} else {
		literal = fmt.Sprintf("%.1f", num)
	}
	s.addToken(NUMBER, numStr, literal)
}

func (s *Scanner) isAlpha(c byte) bool {
	return unicode.IsLetter(rune(c)) || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || unicode.IsDigit(rune(c))
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current+1]
}

// format: <TOKEN_TYPE> <LEXEME> <LITERAL>
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
				s.addToken(EQUAL, "=", "null")
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
				for !s.isAtEnd() && !s.nextMatch('\n') {
					s.advance()
				}
			} else {
				s.addToken(SLASH, string(t), "null")
			}
		case '"':
			s.start = s.current
			// as long as cannot find closing quote
			for !s.nextMatch('"') {
				if s.isAtEnd() || s.nextMatch('\n') {
					s.addError("Unterminated string.")
					break
				}
				s.advance()
			}

			// if closing quote matches
			if s.nextMatch('"') {
				s.advance()
				fullStr := string(s.source[s.start+1 : s.current])
				s.addToken(STRING, fmt.Sprintf("\"%s\"", fullStr), fullStr)
			}
		case ' ': //ignore whitespace
		case '\t': //ignore tab
		case '\r': //ignore carriage returns
		case '\n':
			s.addLine()
		default:
			if unicode.IsDigit(rune(t)) {
				s.addNumber()
			} else if s.isAlpha(t) {
				s.start = s.current
				for s.isAlphaNumeric(s.peek()) {
					s.advance()
				}

				keyword, ok := keywords[string(s.source[s.start:s.current+1])]
				if ok {
					s.addToken(keyword, string(s.source[s.start:s.current+1]), "null")
				} else {
					s.addToken(IDENTIFIER, string(s.source[s.start:s.current+1]), "null")
				}
			} else {
				s.addError(fmt.Sprintf("Unexpected character: %s", string(t)))
			}
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
