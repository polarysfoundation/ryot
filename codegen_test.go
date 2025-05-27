package ryot

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateBytecode_SimpleFunction(t *testing.T) {
	program := &Program{
		Contracts: []*ContractDecl{
			{
				Name: "TestContract",
				Funcs: []*FuncDecl{
					{
						Public:     true,
						Name:       "getValue",
						Args:       []Argument{},
						ReturnType: "uint64",
						Body: []Statement{
							&ReturnStatement{
								Expr: &UInt64Literal{Value: 42},
							},
						},
					},
				},
			},
		},
	}

	expectedBytecode := `FUNC getValue 0
PUSH 42
RETURN`

	bytecode := GenerateBytecode(program)
	if strings.TrimSpace(bytecode) != strings.TrimSpace(expectedBytecode) {
		t.Errorf("Bytecode mismatch.\nExpected:\n%s\nGot:\n%s", expectedBytecode, bytecode)
	}

	// Guardar el bytecode en un archivo
	outputFilename := "tests/test_func.ryc"
	err := os.WriteFile(outputFilename, []byte(bytecode), 0644)
	if err != nil {
		t.Fatalf("Failed to write bytecode to file %s: %v", outputFilename, err)
	} else {
		t.Logf("Bytecode for TestGenerateBytecode_FunctionWithArgsAndBinaryExpr saved to %s", outputFilename)
	}
}

func TestGenerateBytecode_FunctionWithArgsAndBinaryExpr(t *testing.T) {
	source := `
	pragma: "0.1.0";
	
	
	class contract MyContract {
		pub func add(a: uint64, b: uint64): uint64 {
			return (a + b);
		}
	}`

	l := NewLexer(source)
	p := NewParser(l)
	program := p.ParseProgram()

	expectedBytecode := `
FUNC add 2
LOAD_ARG 0
LOAD_ARG 1
ADD
RETURN`

	bytecode := GenerateBytecode(program)
	if strings.TrimSpace(bytecode) != strings.TrimSpace(expectedBytecode) {
		t.Errorf("Bytecode mismatch.\nExpected:\n%s\nGot:\n%s", expectedBytecode, bytecode)
	}

	outputFilename := "tests/test_func_2.ryc"
	err := os.WriteFile(outputFilename, []byte(bytecode), 0644)
	if err != nil {
		t.Fatalf("Failed to write bytecode to file %s: %v", outputFilename, err)
	} else {
		t.Logf("Bytecode for TestGenerateBytecode_FunctionWithArgsAndBinaryExpr saved to %s", outputFilename)
	}
}

func TestBytecodeEmitter_EmitExpression(t *testing.T) {
	// Setup currentFunc for resolveArgIndex
	originalCurrentFunc := currentFunc
	currentFunc = &FuncDecl{
		Args: []Argument{
			{Name: "x", Type: "uint64"},
			{Name: "y", Type: "uint64"},
		},
	}
	defer func() { currentFunc = originalCurrentFunc }() // Restore

	tests := []struct {
		name             string
		expr             Expression
		expectedBytecode string
	}{
		{
			name:             "Identifier",
			expr:             &Identifier{Name: "y"},
			expectedBytecode: "LOAD_ARG 1",
		},
		{
			name:             "UInt64Literal",
			expr:             &UInt64Literal{Value: 123},
			expectedBytecode: "PUSH 123",
		},
		{
			name: "BinaryExpr",
			expr: &BinaryExpr{
				Left:     &Identifier{Name: "x"},
				Operator: "+",
				Right:    &UInt64Literal{Value: 5},
			},
			expectedBytecode: "LOAD_ARG 0\nPUSH 5\nADD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emitter := NewBytecodeEmitter()
			emitter.EmitExpression(tt.expr)
			bytecode := strings.Join(emitter.Instructions, "\n")
			if strings.TrimSpace(bytecode) != strings.TrimSpace(tt.expectedBytecode) {
				t.Errorf("Bytecode mismatch for %s.\nExpected:\n%s\nGot:\n%s", tt.name, tt.expectedBytecode, bytecode)
			}
		})
	}
}

func TestBytecodeEmitter_resolveArgIndex_NotFound(t *testing.T) {
	// Setup currentFunc
	originalCurrentFunc := currentFunc
	currentFunc = &FuncDecl{
		Name: "testFunc",
		Args: []Argument{{Name: "arg1", Type: "uint64"}},
	}
	defer func() { currentFunc = originalCurrentFunc }() // Restore

	emitter := NewBytecodeEmitter()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("resolveArgIndex did not panic for unknown argument")
		} else {
			expectedPanicMsg := "arg not found: unknownArg"
			if r.(string) != expectedPanicMsg {
				t.Errorf("Panic message mismatch. Expected: %q, Got: %q", expectedPanicMsg, r.(string))
			}
		}
	}()

	// This should panic
	emitter.resolveArgIndex("unknownArg")
}

func TestGenerateBytecode_EmptyProgram(t *testing.T) {
	program := &Program{}
	expectedBytecode := ""
	bytecode := GenerateBytecode(program)
	if strings.TrimSpace(bytecode) != strings.TrimSpace(expectedBytecode) {
		t.Errorf("Bytecode mismatch for empty program.\nExpected:\n%s\nGot:\n%s", expectedBytecode, bytecode)
	}
}
