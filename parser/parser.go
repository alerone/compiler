package parser

import (
	"compiler/lexer"
	"fmt"
)

type Parser struct {
	lexer     lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
}

func NewParser(lexer lexer.Lexer) Parser {
	p := Parser{lexer: lexer}
	p.NextToken()
	p.NextToken()
	return p
}

// Return true if the current token matches.
func (self *Parser) checkToken(kind lexer.TokenType) bool {
	return self.curToken.Kind == kind
}

// Return true if the Peek Token matches
func (self *Parser) checkPeek(kind lexer.TokenType) bool {
	return self.peekToken.Kind == kind
}

// Try to match the current token. If not, error. Advances the current token.
func (self *Parser) match(kind lexer.TokenType) {
	if !self.checkToken(kind) {
		self.Abort(fmt.Sprintf("Expected %v, got %v", kind, self.curToken.Kind))
	}
	self.NextToken()
}

// Advances the next Token.
func (self *Parser) NextToken() {
	self.curToken = self.peekToken
	self.peekToken = self.lexer.GetToken()
}

// Advances the next Token.
func (self *Parser) Abort(message string) {
	panic(fmt.Sprintf("Error. %v", message))
}

// program ::= {statement}
func (self *Parser) Program() {
	fmt.Println("PROGRAM")

    for self.checkToken(lexer.NEWLINE){
        self.NextToken()
    }

	// Parse all statements in the program until EOF
	for !self.checkToken(lexer.EOF) {
		self.Statement()
	}
}

func (self *Parser) Statement() {
	// Check the first token to see what kind of statement is

	switch {
	// "PRINT" (expresion | string)
	case self.checkToken(lexer.PRINT):
		fmt.Println("STATEMENT-PRINT")
		self.NextToken()

		if self.checkToken(lexer.STRING) {
			self.NextToken()
		} else {
			self.Expression()
		}
	// "IF" comparison "THEN" {statement} "ENDIF"
	case self.checkToken(lexer.IF):
		fmt.Println("STATEMENT-IF")
		self.NextToken()
		self.Comparison()

		self.match(lexer.THEN)
		self.nl()

		for !self.checkToken(lexer.ENDIF) {
			self.Statement()
		}
		self.match(lexer.ENDIF)
	// "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE"
	case self.checkToken(lexer.WHILE):
		fmt.Println("STATEMENT-WHILE")
		self.NextToken()
		self.Comparison()

		self.match(lexer.REPEAT)
		self.nl()

		for !self.checkToken(lexer.ENDWHILE) {
			self.Statement()
		}
		self.match(lexer.ENDWHILE)
		// "LABEL" ident
	case self.checkToken(lexer.LABEL):
		fmt.Println("STATEMENT-LABEL")
		self.NextToken()
		self.match(lexer.IDENT)
		// "GOTO" ident
	case self.checkToken(lexer.GOTO):
		fmt.Println("STATEMENT-GOTO")
		self.NextToken()
		self.match(lexer.IDENT)
		// "LET" ident "=" expression
	case self.checkToken(lexer.LET):
		fmt.Println("STATEMENT-LET")
		self.NextToken()
		self.match(lexer.IDENT)

		self.match(lexer.EQ)
		self.Expression()
		// "INPUT" ident
	case self.checkToken(lexer.INPUT):
		fmt.Println("STATEMENT-INPUT")
		self.NextToken()
		self.match(lexer.IDENT)
	default:
		self.Abort(fmt.Sprintf("Invalid statement at %v (%v)", self.curToken.Text, self.curToken.Kind))

	}

	self.nl()
}

func (self *Parser) Expression() {
	fmt.Println("EXPRESSION")
}

// 
func (self *Parser) Comparison() {
	fmt.Println("COMPARISON")

}

func (self *Parser) Ident() {
	fmt.Println("IDENT")
}

func (self *Parser) nl() {
	fmt.Println("NEWLINE")

	// Require at least one newline.
	self.match(lexer.NEWLINE)
	// But we will allow more newlines possible
	for self.checkToken(lexer.NEWLINE) {
		self.NextToken()
	}
}
