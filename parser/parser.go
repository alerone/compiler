package parser

import (
	"compiler/emitter"
	"compiler/lexer"
	"fmt"
)

type Parser struct {
	lexer          lexer.Lexer
	emitter        *emitter.Emitter
	curToken       lexer.Token
	peekToken      lexer.Token
	labelsDeclared map[string]bool
	labelsGotoed   map[string]int
	symbols        map[string]bool
}

func NewParser(lexer lexer.Lexer, emitter *emitter.Emitter) Parser {
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

    // For each program we need a main function (library allows us to printf and scanf things).
    self.emitter.HeaderLine("#include <stdio.h>")
    self.emitter.HeaderLine("int main(void){")

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
		self.NextToken()

		if self.checkToken(lexer.STRING) {
            // Simple string, so print it.
            self.emitter.EmitLine("printf(\"" + self.curToken.Text + "\\n\");")
			self.NextToken()
		} else {
            // Expect an expression and print the result as a float.
            self.emitter.Emit("printf(\"%" + ".2f\\n\", (float)(")
			self.expression()
            self.emitter.EmitLine("));")
		}
	// "IF" comparison "THEN" {statement} "ENDIF"
	case self.checkToken(lexer.IF):
		self.NextToken()
        self.emitter.Emit("if(")
		self.comparison()

		self.match(lexer.THEN)
		self.nl()
        self.emitter.EmitLine("){")

		for !self.checkToken(lexer.ENDIF) {
			self.statement()
		}
		self.match(lexer.ENDIF)
        self.emitter.EmitLine("}")
	// "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE"
	case self.checkToken(lexer.WHILE):
		self.NextToken()
        self.emitter.Emit("while(")
		self.comparison()

		self.match(lexer.REPEAT)
		self.nl()
        self.emitter.EmitLine("){")

		for !self.checkToken(lexer.ENDWHILE) {
			self.statement()
		}
		self.match(lexer.ENDWHILE)
        self.emitter.EmitLine("}")
		// "LABEL" ident
	case self.checkToken(lexer.LABEL):
		self.NextToken()

		if self.labelsDeclared[self.curToken.Text] {
			self.Abort(fmt.Sprintf("Label already exists: %v", self.curToken.Text))
		}
		self.labelsDeclared[self.curToken.Text] = true
        self.emitter.EmitLine(self.curToken.Text + ":")
		self.match(lexer.IDENT)
		// "GOTO" ident
	case self.checkToken(lexer.GOTO):
		self.NextToken()
		self.labelsGotoed[self.curToken.Text]++
        self.emitter.EmitLine("goto " + self.curToken.Text + ";")
		self.match(lexer.IDENT)
		// "LET" ident "=" expression
	case self.checkToken(lexer.LET):
		self.NextToken()

		// Check if ident exists in symbols table. If not, declare it.
		if !self.symbols[self.curToken.Text] {
			self.symbols[self.curToken.Text] = true
            self.emitter.HeaderLine("float " + self.curToken.Text + ";")
		}

        self.emitter.Emit(self.curToken.Text + " = ")

		self.match(lexer.IDENT)

		self.match(lexer.EQ)
		self.expression()
        self.emitter.EmitLine(";")
		// "INPUT" ident
	case self.checkToken(lexer.INPUT):
		self.NextToken()

		if !self.symbols[self.curToken.Text] {
			self.symbols[self.curToken.Text] = true
            self.emitter.HeaderLine("float " + self.curToken.Text + ";")
		}

        self.emitter.EmitLine("if(0 == scanf(\"%" + "f\", &" + self.curToken.Text + ")) {")
        self.emitter.EmitLine(self.curToken.Text + " = 0;")
        self.emitter.Emit("scanf(\"%")
        self.emitter.EmitLine("*s\");")
        self.emitter.EmitLine("}")

		self.match(lexer.IDENT)
	default:
		self.Abort(fmt.Sprintf("Invalid statement at %v (%v)", self.curToken.Text, self.curToken.Kind))

	}

	self.nl()
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func (self *Parser) comparison() {
	self.expression()
	// Must be one comparison operator and another expression.
	if self.isComparisonOperator() {
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
		self.expression()
	} else {
		self.Abort(fmt.Sprint("Expected comparison operator at: ", self.curToken.Text))
	}

	for self.isComparisonOperator() {
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
		self.expression()
	}

}

// expression ::= term {( "-" | "+" ) term}
func (self *Parser) expression() {
	self.term()
	for self.checkToken(lexer.MINUS) || self.checkToken(lexer.PLUS) {
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
		self.term()
	}
}

// term ::= unary {( "/" | "*" ) unary}
func (self *Parser) term() {
	self.unary()
	for self.checkToken(lexer.SLASH) || self.checkToken(lexer.ASTERISK) {
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
		self.unary()
	}
}

// unary ::= ["+" | "-"] primary
func (self *Parser) unary() {
	if self.checkToken(lexer.PLUS) || self.checkToken(lexer.MINUS) {
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
	}
	self.primary()
}

// primary ::= number | ident
func (self *Parser) primary() {
	switch {
	case self.checkToken(lexer.NUMBER):
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
	case self.checkToken(lexer.IDENT):
		// Ensure the variable already exists!
		if !self.symbols[self.curToken.Text] {
			self.Abort(fmt.Sprintf("Referencing variable before assignment: %v", self.curToken.Text))
		}
        self.emitter.Emit(self.curToken.Text)
		self.NextToken()
	default:
		self.Abort(fmt.Sprintf("Unexpected token at: %v", self.curToken.Text))
	}
}

func (self *Parser) isComparisonOperator() bool {
	return self.checkToken(lexer.GT) || self.checkToken(lexer.GTEQ) || self.checkToken(lexer.LT) || self.checkToken(lexer.LTEQ) || self.checkToken(lexer.EQ) || self.checkToken(lexer.EQEQ) || self.checkToken(lexer.NOTEQ)
}

func (self *Parser) nl() {
	// Require at least one newline.
	self.match(lexer.NEWLINE)
	// But we will allow more newlines possible
	for self.checkToken(lexer.NEWLINE) {
		self.NextToken()
	}
}
