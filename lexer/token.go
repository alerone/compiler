package lexer


type Token struct {
    Text string
    Kind int
}

var TokenType map[string]int = map[string]int{
    "EOF" : -1,
    "NEWLINE" : 0,
    "NUMBER" : 2,
    "STRING": 3,
    // Keywords
    "LABEL": 101,
    "GOTO": 102,
    "PRINT": 103,
    "INPUT": 104,
    "LET": 105,
    "IF": 106,
    "THEN": 107, 
    "ENDIF": 108,
    "WHILE": 109,
    "REPEAT": 110,
    "ENDWHILE": 111,
    // Operators
    "EQ": 201,
    "PLUS": 202,
    "MINUS": 203,
    "ASTERISK": 204,
    "SLASH": 205,
    "EQEQ": 206,
    "NOTEQ": 207,
    "LT": 208,
    "LETQ": 209,
    "GT": 210,
    "GTEQ": 211,
}
