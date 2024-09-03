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

func addToken(tokens []string, tokenType TokenType, lexeme string, literal string) []string {
	return append(tokens, fmt.Sprintf("%s %s %s", tokenType.String(), lexeme, literal))
}

func nextMatch(i int, input []byte, expected byte) bool {
	if i == (len(input) - 1) {
		return false
	}

	if input[i+1] != expected {
		return false
	}

	return true
}

// output: <token_type> <lexeme> <literal>
func Tokenize(input []byte) ([]string, []string) {
	var tokens []string
	var errors []string

	for i := 0; i < len(input); i++ {
		t := input[i]

		switch t {
		case '(':
			tokens = addToken(tokens, LEFT_PAREN, string(t), "null")
		case ')':
			tokens = addToken(tokens, RIGHT_PAREN, string(t), "null")
		case '{':
			tokens = addToken(tokens, LEFT_BRACE, string(t), "null")
		case '}':
			tokens = addToken(tokens, RIGHT_BRACE, string(t), "null")
		case ',':
			tokens = addToken(tokens, COMMA, string(t), "null")
		case '.':
			tokens = addToken(tokens, DOT, string(t), "null")
		case '-':
			tokens = addToken(tokens, MINUS, string(t), "null")
		case '+':
			tokens = addToken(tokens, PLUS, string(t), "null")
		case ';':
			tokens = addToken(tokens, SEMICOLON, string(t), "null")
		case '/':
			tokens = addToken(tokens, SLASH, string(t), "null")
		case '*':
			tokens = addToken(tokens, STAR, string(t), "null")
		case '=':
			if nextMatch(i, input, '=') {
				tokens = addToken(tokens, EQUAL_EQUAL, "==", "null")
				i += 1
			} else {
				tokens = addToken(tokens, EQUAL, string(t), "null")
			}
		case '!':
			if nextMatch(i, input, '=') {
				tokens = addToken(tokens, BANG_EQUAL, "!=", "null")
				i += 1
			} else {
				tokens = addToken(tokens, BANG, string(t), "null")
			}
		case '\n': //ignore line feeds
		case '\r': //ignore carriage returns
		default:
			errors = append(errors, fmt.Sprintf("[line 1] Error: Unexpected character: %s", string(t)))
		}
	}

	tokens = append(tokens, "EOF  null")

	return tokens, errors
}
