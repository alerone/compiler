package lexer

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	source  string
	CurChar rune
	curPos  int
}

func NewLexer(source string) Lexer {
	s := source + "\n"
	curPos := -1
    l := Lexer{source: s, curPos: curPos}
    l.NextChar()
	return l 
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
	if self.curPos+1 >= len(self.source) {
		return '\x00'
	}
	return rune(self.source[self.curPos+1])
}

// Invalid token found, print error message and exit.
func (l *Lexer) Abort(message string) {
	panic(fmt.Sprintf("lexing err. %v", message))
}

// Skip whitespaces except newlines, which we will indicate the end of a statement.
func (self *Lexer) SkipWhitespace() {
	for self.CurChar == ' ' || self.CurChar == '\t' || self.CurChar == '\r' {
		self.NextChar()
	}
}

// Skip comments in code.
func (self *Lexer) SkipComment() {
	if self.CurChar == '#' {
		for self.CurChar != '\n' {
			self.NextChar()
		}
	}
}

// Return the next token.
func (self *Lexer) GetToken() Token {
	self.SkipWhitespace()
	self.SkipComment()
	var token Token
	stringChar := string(self.CurChar)

	switch {
	case self.CurChar == '=':
		if self.Peek() == '=' {
			lastChar := self.CurChar
			self.NextChar()
			token = Token{string([]rune{lastChar, self.CurChar}), EQEQ}
		} else {
			token = Token{stringChar, EQ}
		}
	case self.CurChar == '>':
		if self.Peek() == '=' {
			lastChar := self.CurChar
			self.NextChar()
			token = Token{string([]rune{lastChar, self.CurChar}), GTEQ}
		} else {
			token = Token{stringChar, GT}
		}
	case self.CurChar == '<':
		if self.Peek() == '=' {
			lastChar := self.CurChar
			self.NextChar()
			token = Token{string([]rune{lastChar, self.CurChar}), LTEQ}
		} else {
			token = Token{stringChar, LT}
		}
	case self.CurChar == '!':
		if self.Peek() == '=' {
			lastChar := self.CurChar
			self.NextChar()
			token = Token{string([]rune{lastChar, self.CurChar}), NOTEQ}
		} else {
			self.Abort(fmt.Sprintf("Expected !=, got !%v", self.Peek()))
		}
	case self.CurChar == '+':
		token = Token{stringChar, PLUS}
	case self.CurChar == '-':
		token = Token{stringChar, MINUS}
	case self.CurChar == '*':
		token = Token{stringChar, ASTERISK}
	case self.CurChar == '"':
		// Get characters between quotations
		self.NextChar()
		startPosition := self.curPos

		for self.CurChar != '"' {
			// Dont allow special characters in the string
			if self.CurChar == '\n' || self.CurChar == '\r' || self.CurChar == '\t' || self.CurChar == '\\' || self.CurChar == '%' {
				self.Abort("illegal character in string.")
			}
			self.NextChar()
		}
		tokText := self.source[startPosition:self.curPos]
		token = Token{tokText, STRING}
	case unicode.IsDigit(self.CurChar):
		// Leading character is a digit so this must be a number.
		// Get all consecutive digits and decimal if there is one.
		startPosition := self.curPos
		for unicode.IsDigit(self.Peek()) {
			self.NextChar()
		}
		if self.Peek() == '.' { // Decimal
			self.NextChar()
			if !unicode.IsDigit(self.Peek()) {
				self.Abort("illegal character in number.")
			}
			for unicode.IsDigit(self.Peek()) {
				self.NextChar()
			}
		}
        


		tokText := self.source[startPosition:self.curPos]
        if self.curPos == startPosition {
            tokText = string(self.CurChar)
        }
        fmt.Println(startPosition, self.curPos)
		token = Token{tokText, NUMBER}
	case self.CurChar == '/':
		token = Token{stringChar, SLASH}
	case self.CurChar == '\n':
		token = Token{stringChar, NEWLINE}
	case self.CurChar == '\x00':
		token = Token{string('\x00'), EOF}
	// keywords
	case unicode.IsLetter(self.CurChar):
		// Leading rune is a letter so this must be an identifier or a keyword
		// Get all consecutive alphanumeric chars
		startPosition := self.curPos
		for isAlphanumeric(self.Peek()) {
			self.NextChar()
		}

		tokText := self.source[startPosition : self.curPos+1]
		keyword, isKey := IsKeyword(tokText)
		if isKey {
			token = Token{tokText, keyword}
		} else {
			token = Token{tokText, IDENT}
		}

	default:
		self.Abort(fmt.Sprintf("unknown token: %v", string(self.CurChar)))
	}

	self.NextChar()
	return token

}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
