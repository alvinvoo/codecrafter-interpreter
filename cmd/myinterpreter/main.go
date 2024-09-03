package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	stringSlice, errorSlice := scanner.Tokenize(fileContents)
	output := ""
	if len(errorSlice) > 0 {
		for _, e := range errorSlice {
			output += e + "\n"
		}
	}

	for _, s := range stringSlice {
		output += s + "\n"
	}

	if len(errorSlice) > 0 {
		fmt.Fprint(os.Stderr, output)
		os.Exit(65)
	} else {
		fmt.Fprint(os.Stdout, output)
		os.Exit(0)
	}
}
