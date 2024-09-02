package scanner

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	VAR
	IDENTIFIER
	EQUAL
	STRING
	SEMICOLON
	EOF
)

func (t TokenType) String() string {
	return [...]string{"LEFT_PAREN", "RIGHT_PAREN", "VAR", "IDENTIFIER", "EQUAL", "STRING", "SEMICOLON", "EOF"}[t]
}

// output: <token_type> <lexeme> <literal>
func Tokenize(input []byte) []string {
	var tokens []string
	for len(input) > 0 {
		t := input[0]

		switch t {
		case '(':
			tokens = append(tokens, "LEFT_PAREN ( null")
		case ')':
			tokens = append(tokens, "RIGHT_PAREN ) null")
		}

		input = input[1:]
	}

	tokens = append(tokens, "EOF  null")

	return tokens
}
