package scanner

import "fmt"

type TokenType int

const (
	// Single character tokens
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	SEMICOLON

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
		"SEMICOLON",
		"EQUAL",
		"IDENTIFIER", "STRING", "NUMBER",
		"VAR",
		"EOF",
	}[t]
}

// output: <token_type> <lexeme> <literal>
func Tokenize(input []byte) ([]string, error) {
	var tokens []string
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
		case '\n': //ignore line feeds
		case '\r': //ignore carriage returns
		default:
			return nil, fmt.Errorf("unexpected character: %q", t)
		}

		input = input[1:]
	}

	tokens = append(tokens, "EOF  null")

	return tokens, nil
}
