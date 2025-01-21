# Compiler made with Go for a very tiny language

## Structure

The compiler's structure is a `Lexer`, a `Parser` and an `Emitter`:

- Lexer: It tokenizes the input to get the keywords from the text
- Parser: With a parsing tree, it validates that the structure of the Very Tiny code matches the Grammar.
- Emitter: It converts the input code to the output code (C code) thanks to the parser.

## Grammar

The Grammar of very tiny is:


```ebnf
program     ::= {statement}
statement   ::= "PRINT" (expression | string) nl
              | "IF" comparison "THEN" nl {statement} "ENDIF" nl
              | "WHILE" comparison "REPEAT" nl {statement} "ENDWHILE" nl
              | "LABEL" ident nl
              | "GOTO" ident nl
              | "LET" ident "=" expression nl
              | "INPUT" ident nl
comparison  ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
expression  ::= term {( "-" | "+" ) term}
term        ::= unary {( "/" | "*" ) unary}
unary       ::= ["+" | "-"] primary
primary     ::= number | ident
nl          ::= "\n"+
```

- `[ ]` : zero or none
- `{ }` : zero or more
- `+`   : one or more of whatever is to the left
- `( )` : just for grouping
- `|`   : logical or

Words are either references to other grammar rules or to tokens defined by the `lexer`.
