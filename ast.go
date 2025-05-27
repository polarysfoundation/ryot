package ryot

type Program struct {
	Contracts []*ContractDecl
}

type ContractDecl struct {
	Version string
	Name    string
	Funcs   []*FuncDecl
}

type FuncDecl struct {
	Public     bool
	Name       string
	Args       []Argument
	ReturnType string
	Body       []Statement
}

type Argument struct {
	Name string
	Type string
}

type Statement interface{}

type ReturnStatement struct {
	Expr Expression
}

type Expression interface{}

type BinaryExpr struct {
	Left     Expression
	Operator string
	Right    Expression
}

type Identifier struct {
	Name string
}

type UInt64Literal struct {
	Value uint64
}
