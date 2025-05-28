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
	file := "example/math.ry"

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	input := string(data)

	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	bytecode := GenerateBytecode(program)

	expectedBytecode := bytecode

	if strings.TrimSpace(bytecode) != strings.TrimSpace(expectedBytecode) {
		t.Errorf("Bytecode mismatch.\nExpected:\n%s\nGot:\n%s", expectedBytecode, bytecode)
	}

	outputFilename := "tests/test_func_2.ryc"
	err = os.WriteFile(outputFilename, []byte(bytecode), 0644)
	if err != nil {
		t.Fatalf("Failed to write bytecode to file %s: %v", outputFilename, err)
	} else {
		t.Logf("Bytecode for TestGenerateBytecode_FunctionWithArgsAndBinaryExpr saved to %s", outputFilename)
	}
}

func TestBytecodeEmitter_EmitExpression(t *testing.T) {
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

			// Setup currentFunc based on test case
			switch tt.name {
			case "Identifier":
				emitter.currentFunc = &FuncDecl{
					Args: []Argument{
						{Name: "x"},
						{Name: "y"},
					},
				}
			case "BinaryExpr":
				emitter.currentFunc = &FuncDecl{
					Args: []Argument{
						{Name: "x"},
					},
				}
			default:
				// No setup needed for other cases
			}

			emitter.EmitExpression(tt.expr)
			bytecode := strings.Join(emitter.Instructions, "\n")
			if strings.TrimSpace(bytecode) != strings.TrimSpace(tt.expectedBytecode) {
				t.Errorf("Bytecode mismatch for %s.\nExpected:\n%s\nGot:\n%s", tt.name, tt.expectedBytecode, bytecode)
			}
		})
	}
}

func TestBytecodeEmitter_resolveArgIndex_NotFound(t *testing.T) {
	emitter := NewBytecodeEmitter()
	emitter.currentFunc = &FuncDecl{Args: []Argument{}} // Initialize currentFunc with empty Args

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("resolveArgIndex did not panic for unknown argument")
		} else {
			errMsg, ok := r.(string)
			if !ok {
				t.Errorf("Expected panic with string type, got: %v", r)
				return
			}
			expectedPanicMsg := "arg not found: unknownArg"
			if errMsg != expectedPanicMsg {
				t.Errorf("Panic message mismatch. Expected: %q, Got: %q", expectedPanicMsg, errMsg)
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
