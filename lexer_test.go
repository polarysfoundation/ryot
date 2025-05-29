package ryot

import "testing"

func TestNextToken(t *testing.T) {
	input := `
	class contract MyContract {
		pub func add(a: uint64, b: uint64): uint64 {
			return (a + b)
		}
	}
	`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{CLASS, "class"},
		{CONTRACT, "contract"},
		{IDENT, "MyContract"},
		{LBRACE, "{"},
		{PUB, "pub"},
		{FUNC, "func"},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COLON, ":"},
		{UINT64, "uint64"},
		{COMMA, ","},
		{IDENT, "b"},
		{COLON, ":"},
		{UINT64, "uint64"},
		{RPAREN, ")"},
		{COLON, ":"},
		{UINT64, "uint64"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{LPAREN, "("},
		{IDENT, "a"},
		{PLUS, "+"},
		{IDENT, "b"},
		{RPAREN, ")"},
		{RBRACE, "}"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestSingleTokens(t *testing.T) {
	input := `+ : ( ) { } , 123 foo_bar`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{PLUS, "+"},
		{COLON, ":"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{COMMA, ","},
		{NUMBER, "123"},
		{IDENT, "foo_bar"},
		{EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestEOF(t *testing.T) {
	l := NewLexer("")
	tok := l.NextToken()
	if tok.Type != EOF {
		t.Fatalf("expected EOF, got %s", tok.Type)
	}
}

func TestUnexpectedChar(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	l := NewLexer("$")
	l.NextToken()
}
