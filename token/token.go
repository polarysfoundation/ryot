package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Especiales
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identificadores + literales
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operadores
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="
	BANG     = "!"
	AND      = "&&"
	OR       = "||"
	MOD      = "%"

	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Palabras clave
	CLASS     = "CLASS"
	STRUCT    = "STRUCT"
	ENUM      = "ENUM"
	PRAGMA    = "PRAGMA"
	PUB       = "PUB"
	PRIV      = "PRIV"
	ST        = "ST"
	FUNC      = "FUNC"
	DELETE    = "DELETE"
	RETURN    = "RETURN"
	NEW       = "NEW"
	CONTRACT  = "CONTRACT"
	INTERFACE = "INTERFACE"

	// Types
	UINT64  = "UINT64"
	ADDRESS = "ADDRESS"
	BOOL    = "BOOL"
	BYTE    = "BYTE"
	TUPLE   = "TUPLE"
)

var keywords = map[string]TokenType{
	"class":     CLASS,
	"struct":    STRUCT,
	"enum":      ENUM,
	"pragma":    PRAGMA,
	"pub":       PUB,
	"priv":      PRIV,
	"st":        ST,
	"func":      FUNC,
	"delete":    DELETE,
	"return":    RETURN,
	"new":       NEW,
	"contract":  CONTRACT,
	"interface": INTERFACE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
