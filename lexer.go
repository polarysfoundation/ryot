package ryot

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	NUMBER    = "NUMBER"
	CLASS     = "CLASS"
	CONTRACT  = "CONTRACT"
	FUNC      = "FUNC"
	PUB       = "PUB"
	RETURN    = "RETURN"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","
	PLUS      = "+"
	SEMICOLON = ";"
	PRAGMA    = "pragma"
	STRING    = "STRING"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	EQ        = "=="
	NOT_EQ    = "!="
	LT        = "<"
	GT        = ">"
	LTE       = "<="
	GTE       = ">="
	PERCENT   = "%"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.position]
	l.readChar() // skip closing quote
	return str
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '"':
		return l.newToken(STRING, l.readString())
	case '+':
		return l.newToken(PLUS, "+")
	case '-':
		return l.newToken(MINUS, "-")
	case '*':
		return l.newToken(ASTERISK, "*")
	case '/':
		return l.newToken(SLASH, "/")
	case ':':
		return l.newToken(COLON, ":")
	case '(':
		return l.newToken(LPAREN, "(")
	case ')':
		return l.newToken(RPAREN, ")")
	case '{':
		return l.newToken(LBRACE, "{")
	case '}':
		return l.newToken(RBRACE, "}")
	case ',':
		return l.newToken(COMMA, ",")
	case ';':
		return l.newToken(SEMICOLON, ";")
	case '%':
		return l.newToken(PERCENT, "%")
	case 0:
		return Token{Type: EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return Token{Type: lookupIdent(literal), Literal: literal}
		} else if isDigit(l.ch) {
			return l.newToken(NUMBER, l.readNumber())
		} else {
			panic(fmt.Sprintf("unexpected char: %q", l.ch))
		}
	}
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	tok := Token{Type: tokenType, Literal: literal}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	switch ident {
	case "pragma":
		return PRAGMA
	case "class":
		return CLASS
	case "contract":
		return CONTRACT
	case "func":
		return FUNC
	case "pub":
		return PUB
	case "return":
		return RETURN
	default:
		return IDENT
	}
}
