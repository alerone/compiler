package lexer

type Lexer struct {
	source  string
	CurChar rune
	curPos  int
}

func NewLexer(source string) Lexer {
	s := source + "\n"
	curPos := -1
	return Lexer{source: s, curPos: curPos}
}

// Process the next character.
func (self *Lexer) NextChar() {
    self.curPos++
    if self.curPos >= len(self.source) {
        self.CurChar = '\x00' // EOF
    } else {
        self.CurChar = rune(self.source[self.curPos])
    }
}

// Return the lookahead char.
func (self *Lexer) Peek() rune {
    if self.curPos + 1 >= len(self.source) {
        return '\x00'
    } 
    return rune(self.source[self.curPos + 1])
}

// Invalid token found, print error message and exit.
func (l *Lexer) Abort(message string) {}

// Skip whitespaces except newlines, which we will indicate the end of a statement.
func (l *Lexer) SkipWhitespace() {}

// Skip comments in code.
func (l *Lexer) SkipComment() {}

// Return the next token.
func (l *Lexer) GetToken() {}
