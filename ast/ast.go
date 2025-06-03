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
	Token  token.Token // Token de tipo 'storage'
	Public bool
	Name   string
	Params []Key
	Value  Value
}

func (sd *StorageDeclaration) statementNode()       {}
func (sd *StorageDeclaration) TokenLiteral() string { return sd.Token.Literal }
func (sd *StorageDeclaration) String() string       { return sd.Name }

// ---------- Keys -------------
type Key struct {
	Token token.Token
	Name  string
	Type  string
}

func (k *Key) statementNode()       {}
func (k *Key) TokenLiteral() string { return k.Token.Literal }
func (k *Key) String() string       { return k.Name }

// ---------- Value -----------
type Value struct {
	Token token.Token
	Type  string
}

func (v *Value) statementNode()       {}
func (v *Value) TokenLiteral() string { return v.Token.Literal }
func (v *Value) String() string       { return v.Type }

// ---------- Storage ----------

type StorageStatement struct {
	Token  token.Token // Token de tipo 'storage'
	Name   string
	Params []Identifier
	Value  Expression
}

func (ss *StorageStatement) expressionNode()      {}
func (ss *StorageStatement) statementNode()       {}
func (ss *StorageStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StorageStatement) String() string       { return ss.Name }

type StorageAccessStatement struct {
	Token  token.Token // Token de tipo 'storage'
	Name   string
	Params []Identifier
}

func (sas *StorageAccessStatement) statementNode()       {}
func (sas *StorageAccessStatement) expressionNode()      {}
func (sas *StorageAccessStatement) TokenLiteral() string { return sas.Token.Literal }
func (sas *StorageAccessStatement) String() string {
	var out bytes.Buffer
	out.WriteString(sas.Name)
	out.WriteString("(")
	params := []string{}
	for _, p := range sas.Params {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

// ---------- Función ----------
type FuncParam struct {
	Name string
	Type string
}

type FuncStatement struct {
	Token      token.Token
	Public     bool
	Name       string
	Params     []Key
	ReturnType Value
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
	Name   string
	Params []Identifier
}

func (ds *DeleteStatement) statementNode()       {}
func (ds *DeleteStatement) expressionNode()      {}
func (ds *DeleteStatement) TokenLiteral() string { return ds.Name }
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

type ByteLiteral struct {
	Token token.Token // The token.BYTE token
	Value uint64
}

func (bl *ByteLiteral) String() string {
	return fmt.Sprintf("%d", bl.Value)
}
func (bl *ByteLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *ByteLiteral) expressionNode() {}

type NewStatement struct {
	Token  token.Token
	Name   string
	Params []Identifier
	Value  Expression
}

func (ns *NewStatement) statementNode()       {}
func (ns *NewStatement) TokenLiteral() string { return "new" }
func (ns *NewStatement) String() string {
	var out bytes.Buffer
	out.WriteString("new ")
	out.WriteString(ns.Name)
	out.WriteString("(")
	params := []string{}
	for _, p := range ns.Params {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("): ")
	out.WriteString(ns.Value.String())
	return out.String()
}

type ConstExpression struct {
	Token token.Token // bool, string, uint64, address, hash, byte
	Name  string
	Value Expression
}

func (ce *ConstExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Name)
	out.WriteString(": ")
	out.WriteString(ce.Value.String())
	return out.String()
}

func (ce *ConstExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *ConstExpression) expressionNode() {}

type ArrayLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) expressionNode() {}

type HashLiteral struct {
	Token token.Token // The 'hash' token
	Value string
}

func (hl *HashLiteral) String() string {
	return hl.Value
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) expressionNode() {}

type ErrLiteral struct {
	Token  token.Token // The 'err' token
	Value  Expression
	Return Expression
}

func (el *ErrLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("err ")
	out.WriteString(el.Value.String())
	out.WriteString(" -> ")
	out.WriteString(el.Return.String())
	return out.String()
}

func (el *ErrLiteral) TokenLiteral() string {
	return el.Token.Literal
}

func (el *ErrLiteral) expressionNode() {}

type ErrValue struct {
	Token token.Token // The 'err' token
	Value Expression
}

func (ev *ErrValue) String() string {
	var out bytes.Buffer
	out.WriteString(ev.Value.String())
	return out.String()
}

func (ev *ErrValue) TokenLiteral() string {
	return ev.Token.Literal
}

func (ev *ErrValue) expressionNode() {}

type VariableStatement struct {
	Token  token.Token // The type token (e.g., uint64, string)
	Name   string
	Value  Expression
	Public bool
}

func (vd *VariableStatement) statementNode()       {}
func (vd *VariableStatement) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vd.Token.Literal) // Type
	out.WriteString(" ")
	out.WriteString(vd.Name) // Variable name
	if vd.Value != nil {
		out.WriteString(": ")
		out.WriteString(vd.Value.String()) // Initial value
	}
	return out.String()
}

type VariableStatementNonInitializer struct {
	Token  token.Token // The type token (e.g., uint64, string)
	Name   string
	Public bool
}

func (vd *VariableStatementNonInitializer) statementNode()       {}
func (vd *VariableStatementNonInitializer) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableStatementNonInitializer) String() string {
	var out bytes.Buffer
	out.WriteString(vd.Token.Literal) // Type
	out.WriteString(" ")
	out.WriteString(vd.Name) // Variable name
	return out.String()
}
