package main

import (
	"compiler/lexer"
	"fmt"
)

func main() {
    source := "LET foobar = 123"
    lx := lexer.NewLexer(source)

    for lx.Peek() != '\x00' {
        fmt.Print(string(lx.CurChar))
        lx.NextChar()
    } 
}
