package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

func readFileAndScan(filename string) []scanner.Token {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	sc := scanner.NewScanner(fileContents)
	sc.Tokenize()
	tokens := sc.GetTokens()
	errorSlice := sc.GetErrors()
	if len(errorSlice) > 0 {
		for _, e := range errorSlice {
			fmt.Fprint(os.Stderr, e+"\n")
		}
		os.Exit(65)
	}

	return tokens
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "tokenize" {
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		sc := scanner.NewScanner(fileContents)
		sc.Tokenize()
		stringSlice := sc.GetTokensString()
		errorSlice := sc.GetErrors()
		if len(errorSlice) > 0 {
			for _, e := range errorSlice {
				fmt.Fprint(os.Stderr, e+"\n")
			}
		}

		for _, s := range stringSlice {
			fmt.Println(s)
		}

		if len(errorSlice) > 0 {
			os.Exit(65)
		} else {
			os.Exit(0)
		}
	} else if command == "parse" {
		tokens := readFileAndScan(os.Args[2])
		if len(tokens) == 0 {
			fmt.Println("No tokens found")
			os.Exit(0)
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		fmt.Println(parser.NewAstPrinter().Print(expr))
	} else if command == "evaluate" {
		tokens := readFileAndScan(os.Args[2])
		if len(tokens) == 0 {
			fmt.Println("No tokens found")
			os.Exit(0)
		}

		p := parser.NewParser(tokens)
		expr, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		interpreter := parser.NewInterpreter()
		values := interpreter.Evaluate(expr)

		if values == nil {
			fmt.Println("nil")
		} else {
			fmt.Printf("%v", values)
		}
	} else if command == "parse_test" {
		b := parser.NewBinary(
			parser.NewLiteral(1),
			scanner.NewToken(scanner.PLUS, "+", ""),
			parser.NewLiteral(2),
		)

		c := parser.NewBinary(
			parser.NewUnary(
				scanner.NewToken(scanner.MINUS, "-", "null"),
				parser.NewLiteral(123),
			),
			scanner.NewToken(scanner.STAR, "*", "null"),
			parser.NewGrouping(parser.NewLiteral(45.67)),
		)

		fmt.Println(parser.NewAstPrinter().Print(b))
		fmt.Println(parser.NewAstPrinter().Print(c))
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
