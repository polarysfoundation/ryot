package parser

import (
	"fmt"
	"strconv"

	"github.com/polarysfoundation/ryot/ast"
	"github.com/polarysfoundation/ryot/lexer"
	"github.com/polarysfoundation/ryot/token"
)

var (
	debug = false
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

	if p.cur.Type != token.STRING_LITERAL { // if the next token is not a string
		p.peekError(token.STRING_LITERAL) // add an error message to the errors slice
		return nil                        // return nil
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
		}

	}

	return stmt // return the ClassStatement node
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

	p.expectPeek(token.LBRACE) // expect the next token to be a left brace

	for p.cur.Type != token.RBRACE && p.cur.Type != token.EOF { // loop until the end of the struct
		p.nextToken() // advance the parser to the next token

		if p.cur.Type == token.IDENT {
			field := ast.StructField{Name: p.cur.Literal} // create a new StructField node
			p.expectPeek(token.COLON)                     // expect the next token to be a colon
			p.nextToken()                                 // advance the parser to the next token

			switch p.cur.Type {
			case token.IDENT: // If the current token is an identifier (e.g. a custom type name)
				field.Type = p.cur.Literal               // Set the field type to the literal value of the token
				p.expectPeek(token.SEMICOLON)            // Expect the next token to be a semicolon
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
			case token.UINT64: // If the current token is a UINT64 type
				field.Type = p.cur.Literal               // Set the field type to "uint64"
				p.expectPeek(token.SEMICOLON)            // Expect the next token to be a semicolon
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
			case token.ADDRESS: // If the current token is an ADDRESS type
				field.Type = p.cur.Literal               // Set the field type to "address"
				p.expectPeek(token.SEMICOLON)            // Expect the next token to be a semicolon
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
			case token.BOOL: // If the current token is a BOOL type
				field.Type = p.cur.Literal               // Set the field type to "bool"
				p.expectPeek(token.SEMICOLON)            // Expect the next token to be a semicolon
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
			case token.LBRACKET: // If the current token is a left bracket (indicating an array type)
				p.expectPeek(token.RBRACKET)                    // Expect the next token to be a right bracket
				p.nextToken()                                   // Advance the parser to the token after the right bracket (which should be the array element type)
				field.Type = fmt.Sprintf("[]%s", p.cur.Literal) // Set the field type to "[]<element_type>"
				stmt.Fields = append(stmt.Fields, field)        // Add the field to the struct's Fields slice
				p.expectPeek(token.SEMICOLON)                   // Expect the next token to be a semicolon
			case token.HASH: // If the current token is a HASH type
				field.Type = p.cur.Literal               // Set the field type to "hash"
				p.expectPeek(token.SEMICOLON)            // Expect the next token to be a semicolon
				stmt.Fields = append(stmt.Fields, field) // Add the field to the struct's Fields slice
			case token.STRING: // If the current token is a STRING type
				field.Type = p.cur.Literal
				p.expectPeek(token.SEMICOLON)
				stmt.Fields = append(stmt.Fields, field)
			default: // If the current token is not a recognized type
				p.peekError(token.IDENT) // Add an error message indicating an identifier was expected
				return nil               // Return nil as parsing failed for this struct field
			}
		}
	}

	p.nextToken() // advance the parser to the next token

	return stmt // return the StructStatement node

}

// parseStorage parses a Storage statement and returns an AST StorageDeclaration node
func (p *Parser) parseStorage(public bool) ast.Statement {
	stmt := &ast.StorageDeclaration{Token: p.cur} // create a new StorageDeclaration node, storing the current token (e.g., 'storage')
	p.nextToken()                                 // advance to the next token (should be the storage name)

	stmt.Public = public // set the Public field based on the 'pub' keyword presence

	if p.cur.Type != token.IDENT { // check if the current token is an identifier (for the storage name)
		p.peekError(token.IDENT) // if not, record an error expecting an identifier
		return nil               // and return nil as parsing failed
	}
	stmt.Name = p.cur.Literal // set the Name field of the StorageDeclaration node

	p.expectPeek(token.LPAREN) // expect the next token to be a left parenthesis '(' for parameters

	stmt.Params = []ast.Key{}                                   // initialize the Params slice
	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF { // loop until a right parenthesis ')' or EOF is encountered
		p.nextToken()                  // advance past the left parenthesis
		if p.cur.Type == token.IDENT { // if the current token is an identifier (parameter name)
			key := ast.Key{Token: p.cur} // create a new Key node
			key.Name = p.cur.Literal     // set the parameter name

			p.expectPeek(token.COLON) // expect the next token to be a colon ':' separating name and type
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

	p.expectPeek(token.COLON) // expect the next token to be a colon ':' separating parameters and value type
	p.nextToken()             // advance past the colon to the value type

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

	p.nextToken() // advance past the value type token

	return stmt // return the fully parsed StorageDeclaration node

}

func (p *Parser) parseFunc(public bool) ast.Statement {
	stmt := &ast.FuncStatement{Token: p.cur}
	p.nextToken()

	stmt.Public = public

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = p.cur.Literal

	fmt.Printf("========= PREPARING FUNC %s ========\n", stmt.Name)

	p.expectPeek(token.LPAREN)

	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
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

	if debug {
		fmt.Println(p.cur.Type)
		fmt.Println(p.cur.Literal)
	}

	stmt.Body = []ast.Statement{}
	for p.cur.Type != token.RBRACE && p.cur.Type != token.EOF {
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
	fmt.Println(expr)
	stmt.Expression = expr
	return stmt
}

func (p *Parser) parseDelete() ast.Statement {
	fmt.Printf("Parsing delete statement: %s \n", p.cur.Literal)

	stmt := &ast.DeleteStatement{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.LPAREN)

	p.nextToken()

	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
		if p.cur.Type == token.IDENT {
			stmt.Params = append(stmt.Params, ast.Identifier{Token: p.cur, Value: p.cur.Literal})
			p.nextToken()
		}
	}

	p.expectPeek(token.SEMICOLON)
	if p.peek.Type == token.RBRACE {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseNew() ast.Statement {
	fmt.Printf("Parsing new statement: %s \n", p.cur.Literal)

	stmt := &ast.NewStatement{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = p.cur.Literal

	if !p.expectPeek(token.LPAREN) {
		p.peekError(token.LPAREN)
		return nil
	}

	p.nextToken()

	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
		if p.cur.Type == token.IDENT {
			stmt.Params = append(stmt.Params, ast.Identifier{Token: p.cur, Value: p.cur.Literal})
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
	case token.LPAREN:
		p.nextToken()
		return p.parseExpression()
	case token.INT:
		return p.parseIntegerLiteral()
	case token.UINT64:
		return p.parseConstExpression()
	case token.STRING_LITERAL:
		fmt.Printf("Parsing string literal: %s \n", p.cur.Literal)
		left = p.parseStringLiteral()
	case token.LBRACKET:
		left = p.parseArrayLiteral()
	}

	if p.peek.Type != token.SEMICOLON {
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
		}
	}

	return left
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	stmt := &ast.ArrayLiteral{Token: token.Token{Type: token.ARRAY, Literal: "array"}}
	p.nextToken()
	for p.cur.Type != token.RBRACKET && p.cur.Type != token.EOF {
		stmt.Elements = append(stmt.Elements, p.parseExpression())
		if p.peek.Type == token.COMMA {
			p.nextToken()
		}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseStringLiteral() ast.Expression {
	stmt := &ast.StringLiteral{Token: p.cur}
	stmt.Value = p.cur.Literal
	p.nextToken()
	return stmt
}

func (p *Parser) parseConstExpression() ast.Expression {
	fmt.Printf("Parsing const expression: %s \n", p.cur.Literal)

	stmt := &ast.ConstExpression{Token: p.cur}
	p.nextToken()

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = p.cur.Literal

	fmt.Println(p.cur.Literal)

	if !p.expectPeek(token.COLON) {
		p.peekError(token.COLON)
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peek.Type != token.SEMICOLON {
		p.peekError(token.SEMICOLON)
		return nil
	}

	p.nextToken()

	return stmt

}

func (p *Parser) parseStorageStatement() ast.Expression {
	fmt.Printf("Parsing storage statement: %s \n", p.cur.Literal)

	stmt := &ast.StorageStatement{Token: token.Token{Type: token.STORAGE, Literal: "storage"}}

	if p.cur.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = p.cur.Literal

	p.expectPeek(token.LPAREN)

	for p.cur.Type != token.RPAREN && p.cur.Type != token.EOF {
		p.nextToken()

		if p.cur.Type == token.IDENT {
			param := ast.Identifier{Token: p.cur, Value: p.cur.Literal}
			stmt.Params = append(stmt.Params, param)
			if p.peek.Type == token.COMMA {
				p.nextToken()
			}
		}
	}

	if !p.expectPeek(token.COLON) {
		fmt.Printf("Parsing storage access statement: %s \n", p.cur.Literal)
		access_storage := &ast.StorageAccessStatement{Token: token.Token{Type: token.STORAGE, Literal: "storage"}}
		access_storage.Name = stmt.Name
		access_storage.Params = stmt.Params

		return access_storage
	}
	p.nextToken()

	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	stmt := &ast.IntegerLiteral{Token: p.cur}
	value, _ := strconv.ParseInt(p.cur.Literal, 10, 64) // parse the literal value as an integer
	stmt.Value = uint64(value)
	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	fmt.Printf("Parsing identifier: %s \n", p.cur.Literal)
	return &ast.Identifier{Token: p.cur, Value: p.cur.Literal}
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
	return expr
}
