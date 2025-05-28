package ryot

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// TestParseProgram valida el parseo completo de un contrato con función pública y return con expresión binaria.
func TestParseProgram(t *testing.T) {
	file := "example/math.ry"

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	input := string(data)

	fmt.Println(input)

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}

	contractLen := len(program.Contracts)
	if contractLen != 1 {
		t.Fatalf("Expected 1 contract, got=%d", contractLen)
	}

	contract := program.Contracts[0]
	if contract.Name != "Math" {
		t.Errorf("Expected contract name 'MyContract', got=%s", contract.Name)
	}

	if len(contract.Funcs) != 5 {
		t.Fatalf("Expected 5 function, got=%d", len(contract.Funcs))
	}

	fn := contract.Funcs[0]
	if fn.Name != "add" {
		t.Errorf("Expected function name 'add', got=%s", fn.Name)
	}
	if !fn.Public {
		t.Errorf("Function should be public")
	}
	if fn.ReturnType != "uint64" {
		t.Errorf("Expected return type 'uint64', got=%s", fn.ReturnType)
	}
	if len(fn.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got=%d", len(fn.Args))
	}

	expectedArgs := []Argument{
		{Name: "a", Type: "uint64"},
		{Name: "b", Type: "uint64"},
	}
	for i, arg := range fn.Args {
		if arg != expectedArgs[i] {
			t.Errorf("Argument %d mismatch. Expected %+v, got %+v", i, expectedArgs[i], arg)
		}
	}

	if len(fn.Body) != 1 {
		t.Fatalf("Expected 1 statement in body, got=%d", len(fn.Body))
	}

	retStmt, ok := fn.Body[0].(*ReturnStatement)
	if !ok {
		t.Fatalf("Expected ReturnStatement, got=%T", fn.Body[0])
	}

	binExpr, ok := retStmt.Expr.(*BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr in return, got=%T", retStmt.Expr)
	}
	if binExpr.Operator != "+" {
		t.Errorf("Expected operator '+', got=%s", binExpr.Operator)
	}
	leftIdent, ok := binExpr.Left.(*Identifier)
	if !ok || leftIdent.Name != "a" {
		t.Errorf("Expected left operand 'a', got=%v", binExpr.Left)
	}
	rightIdent, ok := binExpr.Right.(*Identifier)
	if !ok || rightIdent.Name != "b" {
		t.Errorf("Expected right operand 'b', got=%v", binExpr.Right)
	}
}

// TestParseContract valida el parseo simple de un contrato vacío.
func TestParseContract(t *testing.T) {
	input := `class contract TestContract {}`
	l := NewLexer(input)
	p := NewParser(l)

	contract := &ContractDecl{
		Version: "0.1.0",
	}

	contract = p.parseClassContract(contract)
	if contract == nil {
		t.Fatal("parseClassContract() returned nil")
	}
	if contract.Name != "TestContract" {
		t.Errorf("Expected contract name 'TestContract', got=%s", contract.Name)
	}
	if len(contract.Funcs) != 0 {
		t.Errorf("Expected empty Funcs slice, got=%d", len(contract.Funcs))
	}
}

// TestParseFuncDecl valida el parseo de una función pública con argumentos y tipo de retorno.
func TestParseFuncDecl(t *testing.T) {
	input := `pub func myFunc(arg1: type1, arg2: type2): retType {}`
	l := NewLexer(input)
	p := NewParser(l)

	fn := p.parseFuncDecl(true)
	if fn == nil {
		t.Fatal("parseFuncDecl() returned nil")
	}
	if !fn.Public {
		t.Error("Expected function to be public")
	}
	if fn.Name != "myFunc" {
		t.Errorf("Expected function name 'myFunc', got=%s", fn.Name)
	}
	if len(fn.Args) != 2 {
		t.Fatalf("Expected 2 arguments, got=%d", len(fn.Args))
	}
	expectedArgs := []Argument{
		{Name: "arg1", Type: "type1"},
		{Name: "arg2", Type: "type2"},
	}
	for i, arg := range fn.Args {
		if arg.Name != expectedArgs[i].Name {
			t.Errorf("Argument %d name expected '%s', got '%s'", i, expectedArgs[i].Name, arg.Name)
		}
		if arg.Type != expectedArgs[i].Type {
			t.Errorf("Argument %d type expected '%s', got '%s'", i, expectedArgs[i].Type, arg.Type)
		}
	}
	if fn.ReturnType != "retType" {
		t.Errorf("Expected return type 'retType', got=%s", fn.ReturnType)
	}
	if len(fn.Body) != 0 {
		t.Errorf("Expected empty function body, got %d statements", len(fn.Body))
	}
}

// TestParseReturnStatement valida el parseo de return con identificadores, expresiones binarias y literales.
func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectBinary  bool
		leftIdent     string
		rightIdent    string
	}{
		{`return foo;`, "foo", false, "", ""},
		{`return (bar + baz);`, "", true, "bar", "baz"},
		{`return 123;`, "123", false, "", ""},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)

		stmt := p.parseReturnStatement()
		if stmt == nil {
			t.Fatalf("parseReturnStatement() returned nil for input '%s'", tt.input)
		}

		if tt.expectBinary {
			binExpr, ok := stmt.Expr.(*BinaryExpr)
			if !ok {
				t.Fatalf("Expected *BinaryExpr, got %T for input '%s'", stmt.Expr, tt.input)
			}
			leftIdent, ok := binExpr.Left.(*Identifier)
			if !ok {
				t.Fatalf("Expected *Identifier for left operand, got %T for input '%s'", binExpr.Left, tt.input)
			}
			if leftIdent.Name != tt.leftIdent {
				t.Errorf("Left operand name expected '%s', got '%s' for input '%s'", tt.leftIdent, leftIdent.Name, tt.input)
			}
			rightIdent, ok := binExpr.Right.(*Identifier)
			if !ok {
				t.Fatalf("Expected *Identifier for right operand, got %T for input '%s'", binExpr.Right, tt.input)
			}
			if rightIdent.Name != tt.rightIdent {
				t.Errorf("Right operand name expected '%s', got '%s' for input '%s'", tt.rightIdent, rightIdent.Name, tt.input)
			}
			if binExpr.Operator != "+" {
				t.Errorf("Operator expected '+', got '%s' for input '%s'", binExpr.Operator, tt.input)
			}
		} else {
			switch expr := stmt.Expr.(type) {
			case *Identifier:
				if expr.Name != tt.expectedIdent {
					t.Errorf("Identifier name expected '%s', got '%s' for input '%s'", tt.expectedIdent, expr.Name, tt.input)
				}
			case *UInt64Literal:
				expectedUint, err := strconv.ParseUint(tt.expectedIdent, 10, 64)
				if err != nil {
					t.Fatalf("Invalid expected uint value '%s'", tt.expectedIdent)
				}
				if expr.Value != expectedUint {
					t.Errorf("UInt64Literal value expected %d, got %d for input '%s'", expectedUint, expr.Value, tt.input)
				}
			default:
				t.Fatalf("Unexpected expression type %T for input '%s'", stmt.Expr, tt.input)
			}
		}
	}
}

// TestParseExpression valida el parseo de identificadores y literales uint64.
func TestParseExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedType  interface{}
		expectedValue interface{}
	}{
		{"myVar", &Identifier{}, "myVar"},
		{"42", &UInt64Literal{}, uint64(42)},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)

		expr := p.parseExpression()

		switch e := expr.(type) {
		case *Identifier:
			if _, ok := tt.expectedType.(*Identifier); !ok {
				t.Errorf("Expected Identifier, got %T for input '%s'", e, tt.input)
			}
			if e.Name != tt.expectedValue.(string) {
				t.Errorf("Identifier name expected '%s', got '%s' for input '%s'", tt.expectedValue, e.Name, tt.input)
			}
		case *UInt64Literal:
			if _, ok := tt.expectedType.(*UInt64Literal); !ok {
				t.Errorf("Expected UInt64Literal, got %T for input '%s'", e, tt.input)
			}
			if e.Value != tt.expectedValue.(uint64) {
				t.Errorf("UInt64Literal value expected %d, got %d for input '%s'", tt.expectedValue, e.Value, tt.input)
			}
		default:
			t.Fatalf("Unsupported expression type %T for input '%s'", expr, tt.input)
		}
	}
}

// TestExpectPeekError valida que parseClassContract lance panic si el token esperado no aparece.
func TestExpectPeekError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic but did not panic")
		} else {
			expectedMsg := "expected next token to be IDENT, got EOF instead"
			if r.(string) != expectedMsg {
				t.Errorf("Panic message mismatch. Expected: %q, got: %q", expectedMsg, r.(string))
			}
		}
	}()

	input := `class contract` // Token IDENT esperado después de contract, pero no está
	l := NewLexer(input)
	p := NewParser(l)

	contract := &ContractDecl{
		Version: "0.1.0",
	}

	p.parseClassContract(contract) // Debe provocar panic
}
