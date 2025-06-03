// lexer.go
package lexer

import (
	"unicode"

	"github.com/polarysfoundation/ryot/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NULL
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	l.skipComment()

	var tok token.Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!': // <--- AGREGAR ESTE CASO
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '%':
		tok = newToken(token.MOD, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING_LITERAL
		tok.Literal = l.readString()
	case '<': // Añadir para operadores de comparación
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>': // Añadir para operadores de comparación
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '&': // Añadir para AND
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // Handle single '&' as illegal or a new token
		}
	case '|': // Añadir para OR
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // Handle single '|' as illegal or a new token
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tok.Type = token.LookupIdent(literal)
			tok.Literal = literal
			return tok
		} else if isDigit(l.ch) {
			if l.isAddress() {
				tok.Type = token.ADDRESS_LITERAL
				tok.Literal = l.readAddress()
				return tok
			} else if l.isHash() {
				tok.Type = token.HASH_LITERAL
				tok.Literal = l.readHash()
				return tok
			}
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) isAddress() bool {
	// Guardar estado actual para poder retroceder
	currentPos := l.position
	currentReadPos := l.readPosition
	currentCh := l.ch

	defer func() {
		l.position = currentPos
		l.readPosition = currentReadPos
		l.ch = currentCh
	}()

	// Verificar el patrón completo "1cx" + 30 caracteres hex
	if l.ch != '1' {
		return false
	}
	l.readChar()
	if l.ch != 'c' {
		return false
	}
	l.readChar()
	if l.ch != 'x' {
		return false
	}
	l.readChar()

	// Verificar los 30 caracteres hexadecimales restantes
	for i := 0; i < 30; i++ {
		if l.readPosition > len(l.input) {
			return false
		}
		if !isHexDigit(l.ch) {
			return false
		}
		l.readChar()
	}

	return true
}

func (l *Lexer) isHash() bool {
	// Guardar estado actual para poder retroceder
	currentPos := l.position
	currentReadPos := l.readPosition
	currentCh := l.ch

	defer func() {
		l.position = currentPos
		l.readPosition = currentReadPos
		l.ch = currentCh
	}()

	// Verificar el patrón completo "0x" + 64 caracteres hex
	if l.ch != '0' {
		return false
	}
	l.readChar()
	if l.ch != 'x' {
		return false
	}
	l.readChar()

	// Verificar los 64 caracteres hexadecimales restantes
	for i := 0; i < 64; i++ {
		if l.readPosition > len(l.input) {
			return false
		}
		if !isHexDigit(l.ch) {
			return false
		}
		l.readChar()
	}

	return true
}

func (l *Lexer) readHash() string {
	position := l.position
	// Avanzar los 64 caracteres ( 0 + x + 64 hex)
	for i := 0; i < 66; i++ {
		if l.readPosition > len(l.input) {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readAddress() string {
	position := l.position
	// Avanzar los 33 caracteres (1 + c + x + 30 hex)
	for i := 0; i < 33; i++ {
		if l.readPosition > len(l.input) {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // consume 'x'
		l.readChar() // empieza desde primer dígito
		for isHexDigit(l.ch) {
			l.readChar()
		}
		return l.input[position:l.position]
	}
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
