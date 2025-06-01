package parser

import (
	"fmt"
	"testing"

	"github.com/polarysfoundation/ryot/ast"
	"github.com/polarysfoundation/ryot/lexer"
)

func TestParse_Program(t *testing.T) {
	input := `pragma: "1.0.0";`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program)
}

func TestParse_Class(t *testing.T) {
	input := `pragma: "1.0.0";
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
	input := `pragma: "1.0.0";
		class contract TestEnumContract {
			enum TestEnum: {
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

func TestParse_Struct(t *testing.T) {
	input := `pragma: "1.0.0";
		class contract TestStructContract {
			struct StructTest: {
				data1: uint64;
				data2: string;
				data3: bool;
				data4: address;
				data5: hash;
				data6: []uint64;
			}
		}
		`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[0].(*ast.StructStatement).Fields)

}

func TestParse_Storage(t *testing.T) {
	input := `pragma: "1.0.0";
	class contract testStorageContract {
		    pub storage count(id: uint64): uint64;

			pub func register(id: uint64): void{
				new count(id): 0;
			}

			pub func unregister(id: uint64): void{
				delete count(id);
			}

			pub func get(id: uint64): uint64{
				return count(id);
			}

			pub func getAndReset(id: uint64): uint64{
        		uint64 res: count(id);
				delete count(id);
				return res;
    		}
	}
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[0].(*ast.StorageDeclaration).Value)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[1].(*ast.FuncStatement).Body[0].(*ast.NewStatement).Name)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[1].(*ast.FuncStatement).ReturnType.Type)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[2].(*ast.FuncStatement).Body[0].(*ast.DeleteStatement).Params)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[3].(*ast.FuncStatement).Body[0].(*ast.ReturnStatement).Value.(*ast.StorageAccessStatement).Params)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[4].(*ast.FuncStatement).Body[0].(*ast.ExpressionStatement).Expression.(*ast.ConstExpression).Name)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[4].(*ast.FuncStatement).Body[1].(*ast.DeleteStatement).Name)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[4].(*ast.FuncStatement).Body[2].(*ast.ReturnStatement).Value.(*ast.Identifier))
}

func TestParse_FuncWithReturn(t *testing.T) {
	input := `pragma: "1.0.0";
	class contract testStorageContract {
		    pub func add(a: uint64, b: uint64): uint64{
				return a + b;
			}

			pub func addWithParents(a: uint64, b: uint64): uint64{
				return (a + b);
			}

			pub func name(): string{
				return _name();
			}

			priv func _name(): string {
				return "test";
			}

			pub func uint64Array(): []uint64{
				return [1, 2, 3];
			}

	}
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	fmt.Println(len(program.Statements[1].(*ast.ClassStatement).Body))

	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[1].(*ast.FuncStatement).Body[0].(*ast.ReturnStatement).Value)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[2].(*ast.FuncStatement).Body[0].(*ast.ReturnStatement).Value)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[3].(*ast.FuncStatement).Body[0].(*ast.ReturnStatement).Value)
	fmt.Println(program.Statements[1].(*ast.ClassStatement).Body[4].(*ast.FuncStatement).Body[0].(*ast.ReturnStatement).Value)
}
