package lexer

type Token struct {
	Text string
	Kind TokenType
}

func IsKeyword(tokText string) (TokenType, bool) {
	if kind, exists := keywords[tokText]; exists {
		return kind, true
	}
	return -2, false
}

type TokenType int

const (
	EOF     TokenType = -1
	NEWLINE           = iota
	NUMBER
	IDENT
	STRING
)

const (
	LABEL TokenType = 101 + iota
	GOTO
	PRINT
	INPUT
	LET
	IF
	THEN
	ENDIF
	WHILE
	REPEAT
	ENDWHILE
)

const (
	EQ TokenType = 201 + iota
	PLUS
	MINUS
	ASTERISK
	SLASH
	EQEQ
	NOTEQ
	LT
	LTEQ
	GT
	GTEQ
)

var keywords = map[string]TokenType{
	"LABEL":    LABEL,
	"GOTO":     GOTO,
	"PRINT":    PRINT,
	"INPUT":    INPUT,
	"LET":      LET,
	"IF":       IF,
	"THEN":     THEN,
	"ENDIF":    ENDIF,
	"WHILE":    WHILE,
	"REPEAT":   REPEAT,
	"ENDWHILE": ENDWHILE,
}

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case NEWLINE:
		return "NEWLINE"
	case NUMBER:
		return "NUMBER"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case LABEL:
		return "LABEL"
	case GOTO:
		return "GOTO"
	case PRINT:
		return "PRINT"
	case INPUT:
		return "INPUT"
	case LET:
		return "LET"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ENDIF:
		return "ENDIF"
	case WHILE:
		return "WHILE"
	case REPEAT:
		return "REPEAT"
	case ENDWHILE:
		return "ENDWHILE"
	case EQ:
		return "EQ"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case EQEQ:
		return "EQEQ"
	case NOTEQ:
		return "NOTEQ"
	case LT:
		return "LT"
	case LTEQ:
		return "LTEQ"
	case GT:
		return "GT"
	case GTEQ:
		return "GTEQ"
	default:
		return "UNKNOWN"
	}
}
