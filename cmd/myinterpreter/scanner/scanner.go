package scanner

import "fmt"

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
	EQUAL

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
		"EQUAL",
		"IDENTIFIER", "STRING", "NUMBER",
		"VAR",
		"EOF",
	}[t]
}

// output: <token_type> <lexeme> <literal>
func Tokenize(input []byte) ([]string, []string) {
	var tokens []string
	var errors []string
	for len(input) > 0 {
		t := input[0]

		switch t {
		case '(':
			tokens = append(tokens, fmt.Sprintf("%s %s null", LEFT_PAREN.String(), string(t)))
		case ')':
			tokens = append(tokens, fmt.Sprintf("%s %s null", RIGHT_PAREN.String(), string(t)))
		case '{':
			tokens = append(tokens, fmt.Sprintf("%s %s null", LEFT_BRACE.String(), string(t)))
		case '}':
			tokens = append(tokens, fmt.Sprintf("%s %s null", RIGHT_BRACE.String(), string(t)))
		case ',':
			tokens = append(tokens, fmt.Sprintf("%s %s null", COMMA.String(), string(t)))
		case '.':
			tokens = append(tokens, fmt.Sprintf("%s %s null", DOT.String(), string(t)))
		case '-':
			tokens = append(tokens, fmt.Sprintf("%s %s null", MINUS.String(), string(t)))
		case '+':
			tokens = append(tokens, fmt.Sprintf("%s %s null", PLUS.String(), string(t)))
		case ';':
			tokens = append(tokens, fmt.Sprintf("%s %s null", SEMICOLON.String(), string(t)))
		case '/':
			tokens = append(tokens, fmt.Sprintf("%s %s null", SLASH.String(), string(t)))
		case '*':
			tokens = append(tokens, fmt.Sprintf("%s %s null", STAR.String(), string(t)))
		case '\n': //ignore line feeds
		case '\r': //ignore carriage returns
		default:
			errors = append(errors, fmt.Sprintf("[line 1] Error: Unexpected character: %s", string(t)))
		}

		input = input[1:]
	}

	tokens = append(tokens, "EOF  null")

	return tokens, errors
}
