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
	IDENT           = "IDENT"
	INT             = "INT"
	STRING          = "STRING"
	STRING_LITERAL  = "STRING_LITERAL"
	ADDRESS_LITERAL = "ADDRESS_LITERAL"
	BOOL_LITERAL    = "BOOL_LITERAL"
	BYTE_LITERAL    = "BYTE_LITERAL"
	HASH_LITERAL    = "HASH_LITERAL"
	CXID_LITERAL    = "CXID_LITERAL"

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
	STORAGE   = "STORAGE"
	FUNC      = "FUNC"
	DELETE    = "DELETE"
	RETURN    = "RETURN"
	NEW       = "NEW"
	CONTRACT  = "CONTRACT"
	INTERFACE = "INTERFACE"
	VOID      = "VOID"
	CHECK     = "CHECK"
	ERR       = "ERR"

	// Types
	UINT64  = "UINT64"
	ADDRESS = "ADDRESS"
	BOOL    = "BOOL"
	BYTE    = "BYTE"
	HASH    = "HASH"
	ARRAY   = "ARRAY"
	CXID    = "CXID"
)

var keywords = map[string]TokenType{
	"class":     CLASS,
	"struct":    STRUCT,
	"enum":      ENUM,
	"pragma":    PRAGMA,
	"pub":       PUB,
	"priv":      PRIV,
	"storage":   STORAGE,
	"func":      FUNC,
	"delete":    DELETE,
	"return":    RETURN,
	"new":       NEW,
	"contract":  CONTRACT,
	"interface": INTERFACE,
	"void":      VOID,

	// Types
	"uint64":  UINT64,
	"address": ADDRESS,
	"bool":    BOOL,
	"byte":    BYTE,
	"hash":    HASH,
	"string":  STRING,
	"true":    BOOL_LITERAL,
	"false":   BOOL_LITERAL,
	"null":    VOID,
	"cxid":    CXID,

	"check": CHECK,
	"err":   ERR,

	"==": EQ,
	"!=": NOT_EQ,
	"<=": LTE,
	">=": GTE,
	"&&": AND,
	"||": OR,
}

func (t TokenType) String() string {
	return string(t)
}

func (t Token) String() string {
	return "Token(" + t.Type.String() + ", " + t.Literal + ")"
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
