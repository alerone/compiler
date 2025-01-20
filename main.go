package main

import (
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Very Tiny Compiler")

	if len(os.Args) != 2 {
		panic("Compiler needs source file as argument")
	}

	source, err := os.ReadFile(os.Args[1])
	check(err)

	lxer := lexer.NewLexer(string(source))
	pser := parser.NewParser(lxer)

	pser.Program()
	fmt.Println("Parsing completed.")

}
