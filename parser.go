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
	p.nextToken() // skip COLON

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
			funcDecl := p.parseFuncDecl(true)
			contract.Funcs = append(contract.Funcs, funcDecl)
		}
		p.nextToken()
	}

	return contract
}

func (p *Parser) parseFuncDecl(isPublic bool) *FuncDecl {
	p.expectPeek(FUNC)
	p.expectPeek(IDENT)
	funcName := p.curToken.Literal

	p.expectPeek(LPAREN)

	args := []Argument{}
	for p.peekToken.Type != RPAREN {
		p.nextToken()
		argName := p.curToken.Literal
		p.expectPeek(COLON)
		p.expectPeek(IDENT)
		argType := p.curToken.Literal
		args = append(args, Argument{Name: argName, Type: argType})

		if p.peekToken.Type == COMMA {
			p.nextToken() // skip comma
		}
	}

	p.expectPeek(RPAREN)
	p.expectPeek(COLON)
	p.expectPeek(IDENT)
	returnType := p.curToken.Literal

	p.expectPeek(LBRACE)

	body := []Statement{}
	for p.peekToken.Type != RBRACE {
		p.nextToken()
		if p.curToken.Type == RETURN {
			stmt := p.parseReturnStatement()
			body = append(body, stmt)
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

	for p.peekToken.Type == PLUS {
		p.nextToken() // consume '+'
		p.nextToken() // move to right-hand expression
		right := p.parsePrimary()
		left = &BinaryExpr{
			Left:     left,
			Operator: "+",
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expression {
	switch p.curToken.Type {
	case IDENT:
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

// Utilidad para esperar un token específico sin avanzar el token principal de forma insegura
func (p *Parser) expectPeekReturn(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}
