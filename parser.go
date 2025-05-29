package ryot

import (
	"fmt"
	"strconv"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t TokenType) {
	if p.peekToken.Type != t {
		panic(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	}
	p.nextToken()
}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{}

	for p.curToken.Type != EOF {
		if p.curToken.Type == PRAGMA {
			p.nextToken() // skip 'pragma'
			contract := p.parseVersion()
			if p.curToken.Type == CLASS {
				switch p.peekToken.Type {
				case CONTRACT:
					contract = p.parseClassContract(contract)
					prog.Contracts = append(prog.Contracts, contract)
				default:
					p.errors = append(p.errors, fmt.Sprintf("unsupported class type: %s", p.peekToken.Type))
				}
			}
		} else {
			p.errors = append(p.errors, fmt.Errorf("pragma expected, got %s", p.curToken.Type).Error())
		}
		p.nextToken()
	}

	return prog
}

func (p *Parser) parseVersion() *ContractDecl {
	if p.curToken.Type == COLON {
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected ':' after 'pragma', got %s", p.curToken.Type))
	}

	contract := &ContractDecl{
		Version: p.curToken.Literal,
	}

	p.nextToken() // skip SEMICOLON

	return contract
}

func (p *Parser) parseClassContract(contract *ContractDecl) *ContractDecl {
	p.expectPeek(CONTRACT)
	p.expectPeek(IDENT)
	contractName := p.curToken.Literal
	p.expectPeek(LBRACE)

	contract.Name = contractName

	for p.curToken.Type != RBRACE && p.curToken.Type != EOF {
		if p.curToken.Type == PUB {
			switch p.peekToken.Type {
			case FUNC:
				funcDecl := p.parseFuncDecl(true)
				contract.Funcs = append(contract.Funcs, funcDecl)
			case STORAGE:
				storageDecl := p.parseStorageDecl()
				contract.Storages = append(contract.Storages, storageDecl)
			}

		}
		p.nextToken()
	}

	return contract
}

func (p *Parser) parseStorageDecl() *StorageDecl {
	p.expectPeek(STORAGE)
	p.expectPeek(IDENT)
	storageName := p.curToken.Literal

	p.expectPeek(LPAREN)

	var keyType *Type
	for p.peekToken.Type != RPAREN && p.peekToken.Type != EOF {
		p.nextToken()

		keyType = p.parseTypes()

		if p.peekToken.Type == COMMA {
			p.nextToken() // skip comma
		}
	}

	p.expectPeek(RPAREN)
	p.expectPeek(COLON)

	p.nextToken() // skip COLON
	valueType := p.curToken.Literal

	p.expectPeek(SEMICOLON)

	return &StorageDecl{
		Name:      storageName,
		KeyType:   keyType,
		ValueType: valueType,
	}

}

func (p *Parser) parseFuncDecl(isPublic bool) *FuncDecl {
	p.expectPeek(FUNC)
	p.expectPeek(IDENT)
	funcName := p.curToken.Literal

	p.expectPeek(LPAREN)

	args := []Argument{}
	for p.peekToken.Type != RPAREN && p.peekToken.Type != EOF {
		p.nextToken()
		argName := p.curToken.Literal
		p.expectPeek(COLON)
		p.nextToken()
		argType := p.curToken.Literal
		args = append(args, Argument{Name: argName, Type: argType})

		if p.peekToken.Type == COMMA {
			p.nextToken() // skip comma
		} else if p.peekToken.Type == RPAREN {
			break
		}
	}

	p.expectPeek(RPAREN)
	p.expectPeek(COLON)

	returnType := ""
	if p.peekToken.Type != RBRACE {
		p.nextToken()
		returnType = p.curToken.Literal
	}

	body := []Statement{}
	for p.peekToken.Type != RBRACE {
		p.nextToken()
		switch p.curToken.Type {
		case RETURN:
			stmt := p.parseReturnStatement()
			body = append(body, stmt)
		case NEW:
			stmt := p.parseStorageAssign()
			body = append(body, stmt)
		case UINT64:
			variables := p.parseVariables()
			body = append(body, variables)
		default:
			p.errors = append(p.errors, fmt.Sprintf("unsupported statement type: %s", p.curToken.Type))
		}
	}

	p.expectPeek(RBRACE)

	return &FuncDecl{
		Public:     isPublic,
		Name:       funcName,
		Args:       args,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *Parser) parseVariables() *Variable {
	varType := p.parseTypes()

	p.expectPeek(IDENT)
	varName := p.curToken.Literal
	p.expectPeek(COLON)
	p.nextToken()

	var value Expression
	switch p.curToken.Type {
	case UINT64:
		val, _ := strconv.ParseUint(p.curToken.Literal, 10, 64)
		value = &UInt64Literal{Value: val}
	case IDENT:
		varName := p.curToken.Literal
		switch p.peekToken.Type {
		case LPAREN:
			p.nextToken()
			value = &StorageAccess{
				Var: varName,
				Key: p.parseExpression(),
			}
		}
	}

	p.expectPeek(SEMICOLON)

	return &Variable{
		Type:  varType,
		Name:  varName,
		Value: value,
	}

}

func (p *Parser) parseTypes() *Type {

	switch p.curToken.Type {
	case UINT64:
		return &Type{Name: UINT64}
	case STRING:
		return &Type{Name: STRING}
	case BOOL:
		return &Type{Name: BOOL}
	case ADDRESS:
		return &Type{Name: ADDRESS}
	case IDENT:
		return &Type{Name: p.curToken.Type}
	default:
		panic(fmt.Sprintf("unsupported type: %s", p.curToken.Type))
	}
}

func (p *Parser) parseStorageAssign() *StorageAssign {
	p.expectPeek(IDENT)
	variable := p.curToken.Literal

	p.expectPeek(LPAREN)
	p.expectPeek(IDENT)

	keyValue := &Identifier{Name: p.curToken.Literal}

	p.expectPeek(RPAREN)
	p.expectPeek(COLON)

	if p.peekToken.Type == IDENT {
		p.nextToken()
	}

	var value Expression
	switch p.curToken.Type {
	case UINT64:
		val, _ := strconv.ParseUint(p.curToken.Literal, 10, 64)
		value = &UInt64Literal{Value: val}
	case IDENT:
		varName := p.curToken.Literal
		switch p.peekToken.Type {
		case LPAREN:
			p.nextToken()
			value = &StorageAccess{
				Var: varName,
				Key: p.parseExpression(),
			}

		default:
			value = &Identifier{Name: varName}
		}
	}

	p.expectPeek(SEMICOLON)

	return &StorageAssign{
		Var:   variable,
		Key:   keyValue,
		Value: value,
	}

}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{}

	p.nextToken() // saltar 'return'

	stmt.Expr = p.parseExpression()
	if stmt.Expr == nil {
		p.errors = append(p.errors, "invalid expression in return statement")
		return nil
	}

	if p.peekToken.Type == SEMICOLON {
		p.nextToken() // consumir ';'
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected ';' at the end of return statement, got %s", p.peekToken.Type))
	}

	return stmt
}

func (p *Parser) parseExpression() Expression {
	left := p.parsePrimary()

	switch p.peekToken.Type {
	case PLUS:
		p.nextToken() // consume '+'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "+",
			Right:    right,
		}
	case MINUS:
		p.nextToken() // consume '-'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "-",
			Right:    right,
		}
	case ASTERISK:
		p.nextToken() // consume '-'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "*",
			Right:    right,
		}
	case SLASH:
		p.nextToken() // consume '-'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "/",
			Right:    right,
		}
	case PERCENT:
		p.nextToken() // consume '-'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "%",
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	switch p.curToken.Type {
	case IDENT:
		if p.peekToken.Type == LPAREN {
			return p.parseStorageAccess()
		}
		return &Identifier{Name: p.curToken.Literal}
	case NUMBER:
		val, _ := strconv.ParseUint(p.curToken.Literal, 10, 64)
		return &UInt64Literal{Value: val}
	case LPAREN:
		p.nextToken()
		exp := p.parseExpression()
		if !p.expectPeekReturn(RPAREN) {
			panic("expected ')' after grouped expression")
		}
		return exp
	default:
		panic(fmt.Sprintf("unsupported expression type: %s", p.curToken.Type))
	}
}

func (p *Parser) parseStorageAccess() *StorageAccess {
	name := p.curToken.Literal
	p.nextToken() // Skip name
	p.expectPeek(LPAREN)
	p.nextToken()
	key := p.parseExpression()
	p.expectPeek(RPAREN)

	return &StorageAccess{
		Var: name,
		Key: key,
	}
}

// Utilidad para esperar un token específico sin avanzar el token principal de forma insegura
func (p *Parser) expectPeekReturn(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}
