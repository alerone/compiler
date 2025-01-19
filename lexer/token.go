package lexer


type Token struct {
    text string
    kind int
}

var TokenType map[string]int = map[string]int{
    "eof" : -1,
    "newline" : 0,
    "number" : 2,
    "string": 3,
    // Keywords
    "label": 101,
    "goto": 102,
    "print": 103,
    "input": 104,
    "let": 105,
    "if": 106,
    "then": 107, 
    "endif": 108,
    "while": 109,
    "repeat": 110,
    "endwhile": 111,
    // Operators

}
