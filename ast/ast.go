package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/polarysfoundation/ryot/token"
)

// Nodo raíz
type Node interface {
	TokenLiteral() string
	String() string
}

// Declaración general
type Statement interface {
	Node
	statementNode()
}

// Expresiones (asignaciones, llamadas, etc)
type Expression interface {
	Node
	expressionNode()
}

// ----------------
// Nodos concretos
// ----------------

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String implements the ast.Node interface for Program.
func (p *Program) String() string {
	return "Program"
}

// ---------- Pragma ----------

type PragmaStatement struct {
	Token token.Token
	Value string
}

func (ps *PragmaStatement) statementNode()       {}
func (ps *PragmaStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PragmaStatement) String() string       { return "pragma: \"" + ps.Value + "\";" }

// ---------- Clase ----------

type ClassStatement struct {
	Token       token.Token
	Name        string
	IsInterface bool
	Body        []Statement
}

func (cs *ClassStatement) statementNode()       {}
func (cs *ClassStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ClassStatement) String() string       { return "class " + cs.Name }

// ---------- Enum ----------

type EnumStatement struct {
	Token  token.Token
	Name   string
	Values []string
}

func (es *EnumStatement) statementNode()       {}
func (es *EnumStatement) TokenLiteral() string { return es.Token.Literal }
func (es *EnumStatement) String() string       { return "enum " + es.Name }

// ---------- Struct ----------

type StructField struct {
	Name string
	Type string
}

type StructStatement struct {
	Token  token.Token
	Name   string
	Fields []StructField
}

func (ss *StructStatement) statementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StructStatement) String() string       { return "struct " + ss.Name }

// ---------- Storage ----------

type StorageDeclaration struct {
	Token  token.Token
	Public bool
	Name   string
	Param  string
	Type   string
}

func (sd *StorageDeclaration) statementNode()       {}
func (sd *StorageDeclaration) TokenLiteral() string { return sd.Token.Literal }
func (sd *StorageDeclaration) String() string       { return sd.Name }

// ---------- Función ----------

type FuncParam struct {
	Name string
	Type string
}

type FuncStatement struct {
	Token      token.Token
	Public     bool
	Name       string
	Params     []FuncParam
	ReturnType string
	Body       []Statement
}

func (fs *FuncStatement) statementNode()       {}
func (fs *FuncStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FuncStatement) String() string       { return fs.Name }

// ---------- Otras declaraciones dentro del cuerpo ----------

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string       { return "return ..." }

type DeleteStatement struct {
	Token  token.Token
	Target Expression
}

func (ds *DeleteStatement) statementNode()       {}
func (ds *DeleteStatement) expressionNode()      {}
func (ds *DeleteStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeleteStatement) String() string       { return "delete ..." }

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string       { return es.Expression.String() }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// (Assuming this file defines NewStorageStatement)
type NewStorageStatement struct {
	Token      token.Token
	Identifier *Identifier
	Args       []Expression
	Value      Expression
}

// String implements the ast.Expression interface for NewStorageStatement.
func (nss *NewStorageStatement) String() string {
	args := []string{}
	for _, arg := range nss.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("new %s(%s): %s", nss.Identifier.String(), strings.Join(args, ", "), nss.Value.String())
}

func (nss *NewStorageStatement) TokenLiteral() string {
	return nss.Token.Literal
}

func (n *NewStorageStatement) expressionNode() {}

// BinaryExpression represents a binary operation (e.g., 1 + 2)
type BinaryExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	if be.Left != nil {
		out.WriteString(be.Left.String())
	}
	out.WriteString(" " + be.Operator + " ")
	if be.Right != nil {
		out.WriteString(be.Right.String())
	}
	out.WriteString(")")
	return out.String()
}

func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BinaryExpression) expressionNode() {}

// CallExpression represents a function or storage call expression.
type CallExpression struct {
	Token     token.Token // The token.IDENT token of the function or storage
	Function  Expression  // Identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// IntegerLiteral represents an integer literal in the AST.
type IntegerLiteral struct {
	Token token.Token // The token.INT token
	Value uint64
}

func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() {}

type StringLiteral struct {
	Token token.Token // The token.STRING token
	Value string
}

func (sl *StringLiteral) String() string {
	return sl.Value
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Token token.Token // The token.BOOL token
	Value bool
}

func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) expressionNode() {}

type AddressExpression struct {
	Token token.Token // The token.ADDRESS token
	Value string
}

func (ae *AddressExpression) String() string {
	return ae.Value
}

func (ae *AddressExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AddressExpression) expressionNode() {}
