package ryot

import (
	"fmt"
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	source := `
    class contract Math {
        pub func add(a: uint64, b: uint64): uint64 {
            return (a + b)
        }
    }`

	lexer := NewLexer(source)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	bin, err := CompileToBinary(program)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("tests/math.rybc", bin, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bytecode binario guardado en math.ryotbc")

}
