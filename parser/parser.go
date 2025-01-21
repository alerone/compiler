package parser

import (
	"compiler/emitter"
	"compiler/lexer"
	"fmt"
)

type Parser struct {
	lexer          lexer.Lexer
	emitter        emitter.Emitter
	curToken       lexer.Token
	peekToken      lexer.Token
	labelsDeclared map[string]bool
	labelsGotoed   map[string]int
	symbols        map[string]bool
}

func NewParser(lexer lexer.Lexer, emitter emitter.Emitter) Parser {
	p := Parser{lexer: lexer}

    p.emitter = emitter

	p.labelsDeclared = make(map[string]bool) // Keep track of all labels declared.
	p.labelsGotoed = make(map[string]int) // All labels goto'ed, so we know if they exist or not.
	p.symbols = make(map[string]bool) // Variables we have declared so far.

	p.NextToken()
	p.NextToken() // Call it twice to initialize the current and the peek.
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

// Aborts the compiling with an error message.
func (self *Parser) Abort(message string) {
	panic(fmt.Sprintf("Error. %v", message))
}

// program ::= {statement}
func (self *Parser) Program() {
	fmt.Println("PROGRAM")

    // For each program we need a main function (library allows us to printf and scanf things).
    self.emitter.HeaderLine("#include <stdio.h>")
    self.emitter.HeaderLine("int main(void){")

    fmt.Println("header:",self.emitter.Header)

	for self.checkToken(lexer.NEWLINE) {
		self.NextToken()
	}

	// Parse all statements in the program until EOF
	for !self.checkToken(lexer.EOF) {
		self.statement()
	}

    self.emitter.EmitLine("return 0;")
    self.emitter.EmitLine("}")

	// Check each label referenced in a GOTO is declared.
	for key := range self.labelsGotoed {
		if !self.labelsDeclared[key] {
			self.Abort(fmt.Sprintf("GOTO to a non-declared label: %v", key))
		}
	}
}

func (self *Parser) statement() {
	// Check the first token to see what kind of statement is
	switch {
	// "PRINT" (expresion | string)
	case self.checkToken(lexer.PRINT):
		fmt.Println("STATEMENT-PRINT")
		self.NextToken()

		if self.checkToken(lexer.STRING) {
			self.NextToken()
		} else {
			self.expression()
		}
	// "IF" comparison "THEN" {statement} "ENDIF"
	case self.checkToken(lexer.IF):
		fmt.Println("STATEMENT-IF")
		self.NextToken()
		self.comparison()

		self.match(lexer.THEN)
		self.nl()

		for !self.checkToken(lexer.ENDIF) {
			self.statement()
		}
		self.match(lexer.ENDIF)
	// "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE"
	case self.checkToken(lexer.WHILE):
		fmt.Println("STATEMENT-WHILE")
		self.NextToken()
		self.comparison()

		self.match(lexer.REPEAT)
		self.nl()

		for !self.checkToken(lexer.ENDWHILE) {
			self.statement()
		}
		self.match(lexer.ENDWHILE)
		// "LABEL" ident
	case self.checkToken(lexer.LABEL):
		fmt.Println("STATEMENT-LABEL")
		self.NextToken()

		if self.labelsDeclared[self.curToken.Text] {
			self.Abort(fmt.Sprintf("Label already exists: %v", self.curToken.Text))
		}
		self.labelsDeclared[self.curToken.Text] = true

		self.match(lexer.IDENT)
		// "GOTO" ident
	case self.checkToken(lexer.GOTO):
		fmt.Println("STATEMENT-GOTO")
		self.NextToken()
		self.labelsGotoed[self.curToken.Text]++
		self.match(lexer.IDENT)
		// "LET" ident "=" expression
	case self.checkToken(lexer.LET):
		fmt.Println("STATEMENT-LET")
		self.NextToken()

		// Check if ident exists in symbols table. If not, declare it.
		if !self.symbols[self.curToken.Text] {
			self.symbols[self.curToken.Text] = true
		}

		self.match(lexer.IDENT)

		self.match(lexer.EQ)
		self.expression()
		// "INPUT" ident
	case self.checkToken(lexer.INPUT):
		fmt.Println("STATEMENT-INPUT")
		self.NextToken()

		if !self.symbols[self.curToken.Text] {
			self.symbols[self.curToken.Text] = true
		}
		self.match(lexer.IDENT)
	default:
		self.Abort(fmt.Sprintf("Invalid statement at %v (%v)", self.curToken.Text, self.curToken.Kind))

	}

	self.nl()
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func (self *Parser) comparison() {
	fmt.Println("COMPARISON")

	self.expression()
	// Must be one comparison operator and another expression.
	if self.isComparisonOperator() {
		self.NextToken()
		self.expression()
	} else {
		self.Abort(fmt.Sprint("Expected comparison operator at: ", self.curToken.Text))
	}

	for self.isComparisonOperator() {
		self.NextToken()
		self.expression()
	}

}

// expression ::= term {( "-" | "+" ) term}
func (self *Parser) expression() {
	fmt.Println("EXPRESSION")

	self.term()
	for self.checkToken(lexer.MINUS) || self.checkToken(lexer.PLUS) {
		self.NextToken()
		self.term()
	}
}

// term ::= unary {( "/" | "*" ) unary}
func (self *Parser) term() {
	fmt.Println("TERM")

	self.unary()
	for self.checkToken(lexer.SLASH) || self.checkToken(lexer.ASTERISK) {
		self.NextToken()
		self.unary()
	}
}

// unary ::= ["+" | "-"] primary
func (self *Parser) unary() {
	fmt.Println("UNARY")
	if self.checkToken(lexer.PLUS) || self.checkToken(lexer.MINUS) {
		self.NextToken()
	}
	self.primary()
}

// primary ::= number | ident
func (self *Parser) primary() {
	fmt.Printf("PRIMARY (%v)\n", self.curToken.Text)
	switch {
	case self.checkToken(lexer.NUMBER):
		self.NextToken()
	case self.checkToken(lexer.IDENT):
		// Ensure the variable already exists!
		if !self.symbols[self.curToken.Text] {
			self.Abort(fmt.Sprintf("Referencing variable before assignment: %v", self.curToken.Text))
		}
		self.NextToken()
	default:
		self.Abort(fmt.Sprintf("Unexpected token at: %v", self.curToken.Text))
	}
}

func (self *Parser) isComparisonOperator() bool {
	return self.checkToken(lexer.GT) || self.checkToken(lexer.GTEQ) || self.checkToken(lexer.LT) || self.checkToken(lexer.LTEQ) || self.checkToken(lexer.EQ) || self.checkToken(lexer.EQEQ) || self.checkToken(lexer.NOTEQ)
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
