package main

import (
	"compiler/emitter"
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

	lexeR := lexer.NewLexer(string(source))
    emitteR := emitter.NewEmitter("out.c")
	parseR := parser.NewParser(lexeR, &emitteR)

	parseR.Program()
    emitteR.WriteFile()
	fmt.Println("Compiling completed.")

}
