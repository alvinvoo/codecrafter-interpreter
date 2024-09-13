package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lox"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/util"
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

		p := lox.NewParser(tokens)
		statements, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		fmt.Println(lox.NewAstPrinter().Print(statements[0]))
	} else if command == "evaluate" {
		tokens := readFileAndScan(os.Args[2])
		if len(tokens) == 0 {
			fmt.Println("No tokens found")
			os.Exit(0)
		}

		p := lox.NewParser(tokens)
		statements, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		interpreter := lox.NewInterpreter()

		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "%v\n", r.(util.RuntimeError).Error())
				os.Exit(70)
			}
		}()

		values := interpreter.Evaluate(statements[0])

		if values == nil {
			fmt.Println("nil")
		} else {
			fmt.Printf("%v", values)
		}
	} else if command == "run" {
		tokens := readFileAndScan(os.Args[2])
		if len(tokens) == 0 {
			fmt.Println("No tokens found")
			os.Exit(0)
		}

		p := lox.NewParser(tokens)
		statements, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(65)
		}

		interpreter := lox.NewInterpreter()

		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "%v\n", r.(util.RuntimeError).Error())
				os.Exit(70)
			}
		}()

		interpreter.Interpret(statements)
	} else if command == "parse_test" {
		b := lox.NewBinary(
			lox.NewLiteral(1),
			scanner.NewToken(scanner.PLUS, "+", ""),
			lox.NewLiteral(2),
		)

		c := lox.NewBinary(
			lox.NewUnary(
				scanner.NewToken(scanner.MINUS, "-", "null"),
				lox.NewLiteral(123),
			),
			scanner.NewToken(scanner.STAR, "*", "null"),
			lox.NewGrouping(lox.NewLiteral(45.67)),
		)

		fmt.Println(lox.NewAstPrinter().Print(b))
		fmt.Println(lox.NewAstPrinter().Print(c))
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
