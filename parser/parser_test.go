package parser

import (
	"fmt"
	"testing"

	"github.com/polarysfoundation/ryot/ast"
	"github.com/polarysfoundation/ryot/lexer"
)

func TestParse_Program(t *testing.T) {
	input := `pragma: "1.0.0"`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program)
}

func TestParse_Class(t *testing.T) {
	input := `pragma: "1.0.0" 
	class contract TestClass {
	
	}
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program.Statements)

}

func TestParse_Enum(t *testing.T) {
	input := `pragma: "1.0.0" 
	class contract TestClass {
		enum TestEnum {
			data1; 
			data2;
			data3;
		}
	}
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[0].(*ast.EnumStatement))

}
