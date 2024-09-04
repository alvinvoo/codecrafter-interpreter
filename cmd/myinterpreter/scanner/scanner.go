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

func (s *Scanner) addToken(tokenType TokenType) {
	// TODO: can do improvement of `start` and `current` to keep track of the position of the token
	lexeme := string(s.source[s.start : s.current+1])

	literal := "null"
	if tokenType == STRING {
		literal = lexeme[1 : len(lexeme)-1]
	}

	if tokenType == NUMBER {
		num, err := strconv.ParseFloat(lexeme, 64)
		if err != nil {
			s.addError(fmt.Sprintf("Error parsing number: %s", lexeme))
			return
		}
		// Separate the integer and fractional parts
		_, frac := math.Modf(num)

		// If the fractional part is not zero, it's a float
		if frac != 0 {
			literal = strconv.FormatFloat(num, 'f', -1, 64)
		} else {
			literal = fmt.Sprintf("%.1f", num)
		}
	}

	s.tokens = append(s.tokens, fmt.Sprintf("%s %s %s", tokenType.String(), lexeme, literal))

	s.start = s.current
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
	for !s.isAtEnd() && (unicode.IsDigit(rune(s.source[s.current+1])) || s.nextMatch('.')) {
		s.advance()
	}

	s.addToken(NUMBER)
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
			s.addToken(LEFT_PAREN)
		case ')':
			s.addToken(RIGHT_PAREN)
		case '{':
			s.addToken(LEFT_BRACE)
		case '}':
			s.addToken(RIGHT_BRACE)
		case ',':
			s.addToken(COMMA)
		case '.':
			s.addToken(DOT)
		case '-':
			s.addToken(MINUS)
		case '+':
			s.addToken(PLUS)
		case ';':
			s.addToken(SEMICOLON)
		case '*':
			s.addToken(STAR)
		case '=':
			if s.nextMatch('=') {
				s.advance()
				s.addToken(EQUAL_EQUAL)
			} else {
				s.addToken(EQUAL)
			}
		case '!':
			if s.nextMatch('=') {
				s.advance()
				s.addToken(BANG_EQUAL)
			} else {
				s.addToken(BANG)
			}
		case '<':
			if s.nextMatch('=') {
				s.advance()
				s.addToken(LESS_EQUAL)
			} else {
				s.addToken(LESS)
			}
		case '>':
			if s.nextMatch('=') {
				s.advance()
				s.addToken(GREATER_EQUAL)
			} else {
				s.addToken(GREATER)
			}
		case '/':
			if s.nextMatch('/') {
				// it's a comment, ignore the rest of the line
				for !s.isAtEnd() && !s.nextMatch('\n') {
					s.start += 1
					s.current += 1
				}
			} else {
				s.addToken(SLASH)
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
				s.addToken(STRING)
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
				for s.isAlphaNumeric(s.peek()) {
					s.advance()
				}

				keyword, ok := keywords[string(s.source[s.start:s.current+1])]
				if ok {
					s.addToken(keyword)
				} else {
					s.addToken(IDENTIFIER)
				}
			} else {
				s.addError(fmt.Sprintf("Unexpected character: %s", string(t)))
			}
		}

		// s.advance()
		s.start += 1
		s.current += 1
	}

	s.addEOF()
}

func (s *Scanner) GetTokens() []string {
	return s.tokens
}

func (s *Scanner) GetErrors() []string {
	return s.errors
}
