package parser

import (
	"fmt"

	"github.com/polarysfoundation/ryot/ast"
	"github.com/polarysfoundation/ryot/lexer"
	"github.com/polarysfoundation/ryot/token"
)

// Parser is the main structure for parsing tokens into an AST
type Parser struct {
	l      *lexer.Lexer
	errors []string
	peek   token.Token
	cur    token.Token
	line   int
	col    int
}

// New creates a new Parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,                 // lexer instance
		errors: make([]string, 0), // initialize errors slice
		line:   1,                 // initialize line number
		col:    1,                 // initialize column number
	}
	p.nextToken() // read the first token
	p.nextToken() // read the second token
	return p      // return the parser instance
}

// nextToken advances the parser to the next token
func (p *Parser) nextToken() {
	p.cur = p.peek           // set the current token to the previous token
	p.peek = p.l.NextToken() // get the next token from the lexer
}

// expectPeek checks if the next token is of the expected type and advances the parser if it is
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peek.Type == t { // check if the next token is of the expected type
		p.nextToken() // advance the parser to the next token
		return true   // return true if the next token is of the expected type
	} else { // if the next token is not of the expected type
		p.peekError(t) // add an error message to the errors slice
		return false   // return false
	}
}

// peekError adds an error message to the errors slice if the next token is not of the expected type
func (p *Parser) peekError(t token.TokenType) {
	msg := "expected next token to be " + string(t) + ", got " + string(p.peek.Type) // construct the error message
	p.errors = append(p.errors, msg)                                                 // add the error message to the errors slice
}

// Errors returns a slice of error messages collected during parsing
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram parses the entire program and returns an AST Program node
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}              // create a new Program node
	program.Statements = []ast.Statement{} // initialize the Statements slice

	for p.cur.Type != token.EOF { // loop until the end of the program

		switch p.cur.Type {
		case token.PRAGMA:
			stmt := p.parsePragma() // parse a Pragma statement
			program.Statements = append(program.Statements, stmt)
		case token.CLASS:
			stmt := p.parseClass()
			program.Statements = append(program.Statements, stmt)
		default:
			p.nextToken() // advance the parser to the next token
		}

	}

	return program // return the Program node
}

// parsePragma parses a Pragma statement and returns an AST PragmaStatement node
func (p *Parser) parsePragma() ast.Statement {
	stmt := &ast.PragmaStatement{Token: p.cur} // create a new PragmaStatement node

	p.expectPeek(token.COLON) // expect the next token to be a colon
	p.nextToken()             // advance the parser to the next token

	if p.cur.Type != token.STRING { // if the next token is not a string
		p.peekError(token.STRING) // add an error message to the errors slice
		return nil                // return nil
	}

	stmt.Value = p.cur.Literal // set the Value field of the PragmaStatement node

	p.expectPeek(token.SEMICOLON) // expect the next token to be a semicolon
	p.nextToken()                 // advance the parser to the next token

	return stmt // return the PragmaStatement node
}

// parseClass parses a Class statement and returns an AST ClassStatement node
func (p *Parser) parseClass() ast.Statement {
	stmt := &ast.ClassStatement{Token: p.cur} // create a new ClassStatement node
	p.nextToken()                             // advance the parser to the next token

	// Check if is interface or contract
	if p.cur.Type == token.INTERFACE {
		stmt.IsInterface = true
		p.nextToken() // advance the parser to the next token
	} else if p.cur.Type == token.CONTRACT {
		stmt.IsInterface = false
		p.nextToken() // advance the parser to the next token
	}

	// Get contract name
	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT) // add an error message to the errors slice
		return nil               // return nil
	}
	stmt.Name = p.cur.Literal // set the Name field of the ClassStatement node

	p.expectPeek(token.LBRACE) // expect the next token to be a left brace

	stmt.Body = []ast.Statement{}                               // initialize the Body slice
	for p.cur.Type != token.RBRACE && p.cur.Type != token.EOF { // loop until the end of the class
		p.nextToken() // advance the parser to the next token

		fmt.Println(p.cur.Literal)
		fmt.Println(p.cur.Type)
		switch p.cur.Type {
		case token.ENUM:
			enumStmt := p.parseEnum()
			stmt.Body = append(stmt.Body, enumStmt)
		}

	}

	return stmt // return the ClassStatement node
}

func (p *Parser) parseEnum() ast.Statement {
	stmt := &ast.EnumStatement{Token: p.cur} // create a new EnumStatement node
	p.nextToken()                            // advance the parser to the next token

	if p.cur.Type != token.IDENT { // if the next token is not an identifier
		p.peekError(token.IDENT) // add an error message to the errors slice
		return nil               // return nil
	}
	stmt.Name = p.cur.Literal // set the Name field of the EnumStatement node

	p.expectPeek(token.LBRACE) // expect the next token to be a left brace

	for p.cur.Type != token.RBRACE && p.cur.Type != token.EOF { // loop until the end of the enum
		p.nextToken() // advance the parser to the next token
		if p.cur.Type == token.IDENT {
			stmt.Values = append(stmt.Values, p.cur.Literal) // add the identifier to the Values slice of the EnumStatement node
			p.expectPeek(token.SEMICOLON)                    // expect the next token to be a semicolon
		}
	}

	p.nextToken() // advance the parser to the next token

	return stmt // return the EnumStatement node
}
