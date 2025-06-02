package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/polarysfoundation/ryot/lexer"
	"github.com/polarysfoundation/ryot/parser"
)

func TestCompiler(t *testing.T) {
	source := "../example/example.ry"
	input, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()

	b, err := json.Marshal(program)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))

	contract, err := Compile(string(input))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(contract)

}
