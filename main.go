package main

import (
	"compiler/lexer"
	"fmt"
)

func main() {
    source := "+- */"
    lx := lexer.NewLexer(source)
    
    lx.NextChar()
    token := lx.GetToken()
    for token.Kind != lexer.TokenType["EOF"] {
        fmt.Println(token.Kind)
        token = lx.GetToken()
    } 
}
