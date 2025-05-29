package ryot

import "fmt"

type Program struct {
	Contracts []*ContractDecl
}

type ContractDecl struct {
	Version  string
	Name     string
	Storages []*StorageDecl
	Funcs    []*FuncDecl
}

type FuncDecl struct {
	Public     bool
	Name       string
	Args       []Argument
	ReturnType string
	Body       []Statement
}

type StorageDecl struct {
	Name      string
	KeyType   Expression
	ValueType Expression
}

type Variable struct {
	Name  string
	Type  Expression
	Value Expression
}

type Argument struct {
	Name string
	Type string
}

type StorageAssign struct {
	Var   string
	Key   Expression
	Value Expression
}

type StorageDelete struct {
	Var string
	Key Expression
}

type StorageAccess struct {
	Var string
	Key Expression
}

type Type struct {
	Name TokenType
}

type Statement interface{}
type ReturnStatement struct{ Expr Expression }

type Expression interface{}

type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

type DeleteStatement struct {
	Name string
	Key  Expression
}

type Identifier struct{ Name string }

type UInt64Literal struct{ Value uint64 }
type StringLiteral struct{ Value string }
type BoolLiteral struct{ Value bool }
type NullLiteral struct{}
type ArrayLiteral struct{ Values []Expression }
type AddressLiteral struct{ Value string }
type TupleLiteral struct{ Values []Expression }

func (f *FuncDecl) String() string {
	return fmt.Sprintf("func %s(%v)", f.Name, f.Args)
}
