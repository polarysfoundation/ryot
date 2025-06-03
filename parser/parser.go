package parser

import (
	"encoding/json"
	"fmt"
	"strconv"

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
func (p *Parser) ParseProgram() ast.Node {
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

	b, _ := json.Marshal(program)
	fmt.Println(string(b)) // print the program in JSON format for debugging

	return program // return the Program node
}

// parsePragma parses a Pragma statement and returns an AST PragmaStatement node
func (p *Parser) parsePragma() ast.Statement {
	stmt := &ast.PragmaStatement{Token: p.cur} // create a new PragmaStatement node

	p.expectPeek(token.COLON) // expect the next token to be a colon
	p.nextToken()             // advance the parser to the next token

	if p.cur.Type != token.STRING_LITERAL { // if the next token is not a string
		p.peekError(token.STRING_LITERAL) // add an error message to the errors slice
	}

	stmt.Value = p.cur.Literal // set the Value field of the PragmaStatement node

	p.expectPeek(token.SEMICOLON)

	p.nextToken() // advance the parser to the next token

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
	}
	stmt.Name = p.cur.Literal // set the Name field of the ClassStatement node

	p.expectPeek(token.LBRACE)

	stmt.Body = []ast.Statement{}                                 // initialize the Body slice
	for p.peek.Type != token.RBRACE && p.peek.Type != token.EOF { // loop until the end of the class
		p.nextToken() // advance the parser to the next token

		public := false
		if p.cur.Type == token.PUB {
			public = true
			p.nextToken()
		}

		switch p.cur.Type {
		case token.ENUM:
			enumStmt := p.parseEnum()
			stmt.Body = append(stmt.Body, enumStmt)
		case token.STRUCT:
			structStmt := p.parseStruct()
			stmt.Body = append(stmt.Body, structStmt)
		case token.STORAGE:
			storageStmt := p.parseStorage(public)
			stmt.Body = append(stmt.Body, storageStmt)
		case token.FUNC:
			funcStmt := p.parseFunc(public)
			stmt.Body = append(stmt.Body, funcStmt)
		default:
			switch p.cur.Type {
			case token.UINT64:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			case token.ADDRESS:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			case token.BOOL:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			case token.BYTE:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			case token.HASH:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			case token.STRING:
				stmt.Body = append(stmt.Body, p.parseVariables(public))
			}
		}

	}

	return stmt // return the ClassStatement node
}

func (p *Parser) parseVariables(public bool) ast.Statement {
	stmt := &ast.VariableStatement{Token: p.cur}
	p.nextToken()

	stmt.Public = public

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	if p.peek.Type == token.COLON {
		p.nextToken()
		p.nextToken()

		stmt.Value = p.parseExpression()
	} else {
		varStmt := &ast.VariableStatementNonInitializer{Token: stmt.Token}
		varStmt.Name = stmt.Name
		varStmt.Public = public

		p.expectPeek(token.SEMICOLON)

		return varStmt
	}

	return stmt
}

// parseEnum parses an Enum statement and returns an AST EnumStatement node
func (p *Parser) parseEnum() ast.Statement {
	stmt := &ast.EnumStatement{Token: p.cur} // create a new EnumStatement node
	p.nextToken()                            // advance the parser to the next token

	if p.cur.Type != token.IDENT { // if the next token is not an identifier
		p.peekError(token.IDENT) // add an error message to the errors slice
		return nil               // return nil
	}
	stmt.Name = p.cur.Literal // set the Name field of the EnumStatement node

	p.nextToken() // skip to COLON

	if !p.expectPeek(token.LBRACE) {
		p.peekError(token.LBRACE)
		return nil
	}

	for p.peek.Type != token.RBRACE && p.peek.Type != token.EOF { // loop until the end of the enum
		p.nextToken() // advance the parser to the next token
		if p.cur.Type == token.IDENT {
			stmt.Values = append(stmt.Values, p.cur.Literal) // add the identifier to the Values slice of the EnumStatement node
			p.expectPeek(token.SEMICOLON)
		}
	}

	p.nextToken() // advance the parser to the next token

	return stmt // return the EnumStatement node
}

// parseStruct parses a Struct statement and returns an AST StructStatement node
func (p *Parser) parseStruct() ast.Statement {
	stmt := &ast.StructStatement{Token: p.cur} // create a new StructStatement node
	p.nextToken()                              // advance the parser to the next token

	if p.cur.Type != token.IDENT { // if the next token is not an identifier
		p.peekError(token.IDENT) // add an error message to the errors slice
		return nil               // return nil
	}

	stmt.Name = p.cur.Literal // set the Name field of the StructStatement node
	p.nextToken()             // advance the parser to the next token

	if !p.expectPeek(token.LBRACE) {
		p.peekError(token.LBRACE)
		return nil
	}

	for p.peek.Type != token.RBRACE && p.peek.Type != token.EOF { // loop until the end of the struct
		p.nextToken() // advance the parser to the next token

		if p.cur.Type == token.IDENT {
			field := ast.StructField{Name: p.cur.Literal} // create a new StructField node

			if !p.expectPeek(token.COLON) {
				p.peekError(token.COLON)
				return nil
			}

			p.nextToken() // advance the parser to the next token

			switch p.cur.Type {
			case token.IDENT: // If the current token is an identifier (e.g. a custom type name)
				field.Type = p.cur.Literal               // Set the field type to the literal value of the token
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.UINT64: // If the current token is a UINT64 type
				field.Type = p.cur.Literal               // Set the field type to "uint64"
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.ADDRESS: // If the current token is an ADDRESS type
				field.Type = p.cur.Literal               // Set the field type to "address"
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.BOOL: // If the current token is a BOOL type
				field.Type = p.cur.Literal               // Set the field type to "bool"
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.LBRACKET: // If the current token is a left bracket (indicating an array type)
				// Expect the next token to be a right bracket
				if !p.expectPeek(token.RBRACKET) {
					p.peekError(token.RBRACKET)
					return nil
				}

				p.nextToken()                                   // Advance the parser to the token after the right bracket (which should be the array element type)
				field.Type = fmt.Sprintf("[]%s", p.cur.Literal) // Set the field type to "[]<element_type>"
				stmt.Fields = append(stmt.Fields, field)        // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.HASH: // If the current token is a HASH type
				field.Type = p.cur.Literal               // Set the field type to "hash"
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			case token.STRING: // If the current token is a STRING type
				field.Type = p.cur.Literal
				stmt.Fields = append(stmt.Fields, field)
				if !p.expectPeek(token.SEMICOLON) {
					p.peekError(token.SEMICOLON)
					return nil
				}
			default: // If the current token is not a recognized type
				p.peekError(token.IDENT) // Add an error message indicating an identifier was expected
				return nil               // Return nil as parsing failed for this struct field
			}
		}
	}

	p.nextToken()

	return stmt // return the StructStatement node

}

// parseStorage parses a Storage statement and returns an AST StorageDeclaration node
func (p *Parser) parseStorage(public bool) ast.Statement {
	fmt.Printf("========= PREPARING STORAGE %s ========\n", p.cur.Literal)

	stmt := &ast.StorageDeclaration{Token: p.cur} // create a new StorageDeclaration node, storing the current token (e.g., 'storage')
	p.nextToken()                                 // advance to the next token (should be the storage name)

	stmt.Public = public // set the Public field based on the 'pub' keyword presence

	if p.cur.Type != token.IDENT { // check if the current token is an identifier (for the storage name)
		p.peekError(token.IDENT) // if not, record an error expecting an identifier
		return nil               // and return nil as parsing failed
	}
	stmt.Name = p.cur.Literal // set the Name field of the StorageDeclaration node

	p.expectPeek(token.LPAREN)

	stmt.Params = []ast.Key{}                                     // initialize the Params slice
	for p.peek.Type != token.RPAREN && p.peek.Type != token.EOF { // loop until a right parenthesis ')' or EOF is encountered
		p.nextToken()                  // advance past the left parenthesis
		if p.cur.Type == token.IDENT { // if the current token is an identifier (parameter name)
			key := ast.Key{Token: p.cur} // create a new Key node
			key.Name = p.cur.Literal     // set the parameter name

			fmt.Println(key)

			p.expectPeek(token.COLON)

			p.nextToken()

			switch p.cur.Type { // determine the parameter type based on the current token
			case token.UINT64:
				key.Type = "uint64"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.ADDRESS:
				key.Type = "address"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.BOOL:
				key.Type = "bool"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.BYTE:
				key.Type = "byte"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.HASH:
				key.Type = "hash"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.STRING:
				key.Type = "string"
				stmt.Params = append(stmt.Params, key) // add the parsed key to the Params slice
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			default: // if the token is not a recognized type
				p.peekError(token.IDENT) // record an error expecting a type identifier
				return nil               // and return nil as parsing failed
			}
		}
	}

	p.nextToken()

	if !p.expectPeek(token.COLON) {
		p.peekError(token.COLON)
		return nil
	}

	p.nextToken() // advance past the colon to the value type

	stmt.Value = ast.Value{Token: p.cur} // create a new Value node for the storage value type
	switch p.cur.Type {                  // determine the value type based on the current token
	case token.UINT64:
		stmt.Value.Type = "uint64"
	case token.ADDRESS:
		stmt.Value.Type = "address"
	case token.BOOL:
		stmt.Value.Type = "bool"
	case token.BYTE:
		stmt.Value.Type = "byte"
	case token.HASH:
		stmt.Value.Type = "hash"
	default: // if the token is not a recognized type
		p.peekError(token.IDENT) // record an error expecting a type identifier
		return nil               // and return nil as parsing failed
	}

	p.expectPeek(token.SEMICOLON)

	return stmt // return the fully parsed StorageDeclaration node

}

func (p *Parser) parseFunc(public bool) ast.Statement {
	stmt := &ast.FuncStatement{Token: p.cur}
	p.nextToken()

	stmt.Public = public

	if p.cur.Type != token.IDENT && p.peek.Type != token.LPAREN {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	fmt.Printf("========= PREPARING FUNC %s ========\n", stmt.Name)

	p.expectPeek(token.LPAREN)

	for p.peek.Type != token.RPAREN && p.peek.Type != token.EOF {
		p.nextToken()

		if p.cur.Type == token.IDENT {
			key := ast.Key{Token: p.cur}
			key.Name = p.cur.Literal

			p.expectPeek(token.COLON)

			p.nextToken()

			switch p.cur.Type {
			case token.UINT64:
				key.Type = "uint64"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.ADDRESS:
				key.Type = "address"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.BOOL:
				key.Type = "bool"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.BYTE:
				key.Type = "byte"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.HASH:
				key.Type = "hash"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.STRING:
				key.Type = "string"
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			case token.LBRACKET:
				p.expectPeek(token.RBRACKET)
				p.nextToken()
				key.Type = fmt.Sprintf("[]%s", p.cur.Literal)
				stmt.Params = append(stmt.Params, key)
				if p.peek.Type == token.COMMA {
					p.nextToken()
				}
			default:
				p.peekError(token.IDENT)
				return nil
			}
		}
	}

	p.nextToken()

	p.expectPeek(token.COLON)

	p.nextToken()

	switch p.cur.Type {
	case token.UINT64:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "uint64"
	case token.ADDRESS:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "address"
	case token.BOOL:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "bool"
	case token.BYTE:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "byte"
	case token.HASH:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "hash"
	case token.VOID:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "void"
	case token.STRING:
		stmt.ReturnType.Token = p.cur
		stmt.ReturnType.Type = "string"
	case token.LBRACKET:
		p.expectPeek(token.RBRACKET)
		p.nextToken()
		stmt.ReturnType.Token = token.Token{Type: token.ARRAY, Literal: "array"}
		stmt.ReturnType.Type = fmt.Sprintf("[]%s", p.cur.Literal)
	default:
		p.peekError(token.IDENT)
		return nil
	}

	p.expectPeek(token.LBRACE)

	stmt.Body = []ast.Statement{}
	for p.peek.Type != token.RBRACE && p.peek.Type != token.EOF {
		p.nextToken()

		switch p.cur.Type {
		case token.RETURN:
			stmt.Body = append(stmt.Body, p.parseReturn())
		case token.NEW:
			stmt.Body = append(stmt.Body, p.parseNew())
		case token.DELETE:
			stmt.Body = append(stmt.Body, p.parseDelete())
		default:
			stmt.Body = append(stmt.Body, p.parseExpressionStatement())
		}
	}

	p.nextToken()

	fmt.Printf("========= FUNC DONE %s ========\n", stmt.Name)

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.cur}
	expr := p.parseExpression()
	stmt.Expression = expr
	return stmt
}

func (p *Parser) parseDelete() ast.Statement {
	fmt.Printf("Parsing delete statement: %s \n", p.cur.Literal)

	stmt := &ast.DeleteStatement{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.LPAREN)

	for p.peek.Type != token.RPAREN && p.peek.Type != token.EOF {
		p.nextToken()
		if p.cur.Type == token.IDENT {
			stmt.Params = append(stmt.Params, ast.Identifier{Token: p.cur, Value: p.cur.Literal})
			if p.peek.Type == token.COMMA {
				p.nextToken()
			}
		}
	}

	p.nextToken()

	p.expectPeek(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseNew() ast.Statement {
	fmt.Printf("Parsing new statement: %s \n", p.cur.Literal)

	stmt := &ast.NewStatement{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.LPAREN)

	p.nextToken()

	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
		if p.cur.Type == token.IDENT {
			stmt.Params = append(stmt.Params, ast.Identifier{Token: p.cur, Value: p.cur.Literal})
			if p.peek.Type == token.COMMA {
				p.nextToken()
			}
			p.nextToken()
		}
	}

	p.expectPeek(token.COLON)
	p.nextToken()

	stmt.Value = p.parseExpression()

	return stmt

}

func (p *Parser) parseReturn() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.cur}
	p.nextToken()

	fmt.Printf("Parsing return statement: %s \n", p.cur.Literal)

	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	fmt.Printf("Parsing expresion: %s \n", p.cur.Literal)

	var left ast.Expression
	switch p.cur.Type {
	case token.IDENT:
		if p.peek.Type != token.LPAREN {
			left = p.parseIdentifier()
		} else {
			left = p.parseStorageStatement()
		}
	case token.CHECK:
		left = p.parseErrLiteral()
	case token.ERR:
		left = p.parseErrValue()
	case token.LPAREN:
		p.nextToken()
		left = p.parseExpression()
	case token.INT:
		fmt.Printf("Parsing integer literal: %s \n", p.cur.Literal)
		left = p.parseIntegerLiteral()
	case token.UINT64:
		left = p.parseConstExpression()
	case token.STRING_LITERAL:
		fmt.Printf("Parsing string literal: %s \n", p.cur.Literal)
		left = p.parseStringLiteral()
	case token.BOOL_LITERAL:
		fmt.Printf("Parsing bool literal: %s \n", p.cur.Literal)
		left = p.parseBoolLiteral()
	case token.HASH_LITERAL:
		fmt.Printf("Parsing hash literal: %s \n", p.cur.Literal)
		left = p.parseHashLiteral()
	case token.LBRACKET:
		left = p.parseArrayLiteral()
	case token.ADDRESS:
		left = p.parseConstExpression()
	case token.BOOL:
		left = p.parseConstExpression()
	case token.BYTE:
		left = p.parseConstExpression()
	case token.HASH:
		left = p.parseConstExpression()
	case token.STRING:
		left = p.parseConstExpression()
	case token.ADDRESS_LITERAL:
		left = p.parseAddressLiteral()

	}

	if p.cur.Type != token.SEMICOLON {
		p.nextToken()
		switch p.cur.Type {
		case token.PLUS:
			return p.parseBinaryExpression(left)
		case token.MINUS:
			return p.parseBinaryExpression(left)
		case token.ASTERISK:
			return p.parseBinaryExpression(left)
		case token.SLASH:
			return p.parseBinaryExpression(left)
		case token.LT:
			return p.parseBinaryExpression(left)
		case token.GT:
			return p.parseBinaryExpression(left)
		case token.EQ:
			return p.parseBinaryExpression(left)
		case token.NOT_EQ:
			return p.parseBinaryExpression(left)
		case token.LTE:
			return p.parseBinaryExpression(left)
		case token.GTE:
			return p.parseBinaryExpression(left)
		case token.AND:
			return p.parseBinaryExpression(left)
		case token.OR:
			return p.parseBinaryExpression(left)
		case token.MOD:
			return p.parseBinaryExpression(left)
		}
	}

	return left
}

func (p *Parser) parseErrValue() ast.Expression {
	stmt := &ast.ErrValue{Token: p.cur}
	p.expectPeek(token.COLON)
	p.nextToken()

	if p.cur.Type != token.STRING_LITERAL {
		p.peekError(token.STRING_LITERAL)
	}

	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseErrLiteral() ast.Expression {
	stmt := &ast.ErrLiteral{Token: p.cur}
	p.expectPeek(token.LPAREN)
	p.nextToken()

	left := p.parseExpression()

	stmt.Value = left

	if p.cur.Type == token.COMMA {
		p.nextToken()
		value := p.parseExpression()
		stmt.Return = value
	}

	return stmt
}

func (p *Parser) parseHashLiteral() ast.Expression {
	stmt := &ast.HashLiteral{Token: token.Token{Type: token.HASH, Literal: "hash"}}
	stmt.Value = p.cur.Literal

	if p.peek.Type != token.COMMA {
		p.expectPeek(token.SEMICOLON)
	}

	return stmt
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	fmt.Printf("Parsing array literal: %s \n", p.cur.Literal)

	stmt := &ast.ArrayLiteral{Token: token.Token{Type: token.ARRAY, Literal: "array"}}

	stmt.Elements = []ast.Expression{}
	for p.peek.Type != token.RBRACKET && p.peek.Type != token.EOF {
		p.nextToken()

		switch p.cur.Type {
		case token.IDENT:
			stmt.Elements = append(stmt.Elements, p.parseIdentifier())
		case token.INT:
			stmt.Elements = append(stmt.Elements, p.parseIntegerLiteral())
		case token.UINT64:
			stmt.Elements = append(stmt.Elements, p.parseConstExpression())
		case token.STRING_LITERAL:
			stmt.Elements = append(stmt.Elements, p.parseStringLiteral())
		case token.BOOL_LITERAL:
			stmt.Elements = append(stmt.Elements, p.parseBoolLiteral())
		case token.LBRACKET:
			stmt.Elements = append(stmt.Elements, p.parseArrayLiteral())
		case token.ADDRESS_LITERAL:
			stmt.Elements = append(stmt.Elements, p.parseAddressLiteral())
		}
	}
	p.nextToken()

	p.expectPeek(token.SEMICOLON)

	return stmt
}

func (p *Parser) parseBoolLiteral() ast.Expression {
	stmt := &ast.BooleanLiteral{Token: p.cur}
	if p.cur.Literal == "true" {
		stmt.Value = true
	} else {
		stmt.Value = false
	}
	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}
	return stmt
}

func (p *Parser) parseStringLiteral() ast.Expression {
	stmt := &ast.StringLiteral{Token: token.Token{Type: token.STRING, Literal: "string"}}
	stmt.Value = p.cur.Literal
	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}
	return stmt
}

func (p *Parser) parseConstExpression() ast.Expression {
	fmt.Printf("Parsing const expression: %s \n", p.cur.Literal)

	stmt := &ast.ConstExpression{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.COLON)

	p.nextToken()

	stmt.Value = p.parseExpression()

	return stmt

}

func (p *Parser) parseStorageStatement() ast.Expression {
	fmt.Printf("Parsing storage statement: %s \n", p.cur.Literal)

	stmt := &ast.StorageStatement{Token: token.Token{Type: token.STORAGE, Literal: "storage"}}

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.LPAREN)

	for p.peek.Type != token.RPAREN && p.peek.Type != token.EOF {
		p.nextToken()

		if p.cur.Type == token.IDENT {
			param := ast.Identifier{Token: p.cur, Value: p.cur.Literal}
			stmt.Params = append(stmt.Params, param)
			if p.peek.Type == token.COMMA {
				p.nextToken()
			}
		}
	}

	p.nextToken()

	if p.peek.Type != token.COLON {
		fmt.Printf("Parsing storage access statement: %s \n", p.cur.Literal)
		access_storage := &ast.StorageAccessStatement{Token: token.Token{Type: token.STORAGE, Literal: "storage"}}
		access_storage.Name = stmt.Name
		access_storage.Params = stmt.Params

		p.expectPeek(token.SEMICOLON)

		return access_storage
	}
	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	stmt := &ast.IntegerLiteral{Token: p.cur}
	value, _ := strconv.ParseInt(p.cur.Literal, 10, 64) // parse the literal value as an integer
	stmt.Value = uint64(value)

	if p.peek.Type == token.RPAREN {
		p.nextToken()
	}

	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}
	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	fmt.Printf("Parsing identifier: %s \n", p.cur.Literal)
	stmt := &ast.Identifier{Token: p.cur, Value: p.cur.Literal}
	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}

	return stmt
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	fmt.Printf("Parsing binary expression: %s \n", p.cur.Literal)

	expr := &ast.BinaryExpression{
		Token:    p.cur,
		Operator: p.cur.Literal,
		Left:     left,
	}
	p.nextToken()
	expr.Right = p.parseExpression()
	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}

	return expr
}

func (p *Parser) parseAddressLiteral() ast.Expression {
	stmt := &ast.AddressExpression{Token: p.cur}
	stmt.Value = p.cur.Literal
	if p.peek.Type == token.SEMICOLON {
		p.expectPeek(token.SEMICOLON)
	}
	return stmt
}
