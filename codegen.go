package ryot

import (
	"fmt"
	"strings"
)

const (
	STORAGE_DECL    = "DECL_STORAGE" // Declara variable de storage
	STORAGE_LOAD    = "SLOAD"        // Carga valor del storage
	STORAGE_STORE   = "SSTORE"       // Guarda valor en storage
	STORAGE_DELETE  = "SDELETE"      // Elimina del storage
	MSTORAGE_LOAD   = "MLOAD"        // Carga valor de memoria
	MSTORAGE_STORE  = "MSTORE"       // Guarda valor en memoria
	MSTORAGE_DELETE = "MDELETE"      // Elimina valor de memoria

)

type BytecodeEmitter struct {
	Instructions []string
	currentFunc  *FuncDecl
}

func NewBytecodeEmitter() *BytecodeEmitter {
	return &BytecodeEmitter{Instructions: []string{}}
}

func (e *BytecodeEmitter) Emit(instr string, args ...interface{}) {
	line := instr
	if len(args) > 0 {
		line += " " + fmt.Sprint(args...)
	}
	e.Instructions = append(e.Instructions, line)
}

func (e *BytecodeEmitter) EmitFunction(f *FuncDecl) {
	e.currentFunc = f
	e.Emit(fmt.Sprintf("FUNC %s %d", f.Name, len(f.Args)))

	for _, stmt := range f.Body {
		switch s := stmt.(type) {
		case *StorageAssign:
			e.EmitExpression(s.Key)
			e.EmitExpression(s.Value)
			e.Emit(STORAGE_STORE, s.Var)
		case *ReturnStatement:
			e.EmitExpression(s.Expr)
			e.Emit("RETURN")
		case *Variable:
			e.EmitExpression(s.Value)
			e.Emit(MSTORAGE_STORE, s.Name)
			e.Emit(MSTORAGE_LOAD, s.Name)
			e.EmitExpression(s.Value)
		default:
			panic(fmt.Sprintf("unsupported stmt: %T", s))
		}
	}
}

func (e *BytecodeEmitter) EmitExpression(expr Expression) {
	switch exp := expr.(type) {
	case *BinaryExpr:
		e.EmitExpression(exp.Left)
		e.EmitExpression(exp.Right)
		switch exp.Operator {
		case "+":
			e.Emit("ADD")
		case "-":
			e.Emit("SUB")
		case "*":
			e.Emit("MUL")
		case "/":
			e.Emit("DIV")
		case "%":
			e.Emit("MOD")
		case "==":
			e.Emit("EQ")
		case "!=":
			e.Emit("NEQ")
		case "<":
			e.Emit("LT")
		case "<=":
			e.Emit("LTE")
		case ">":
			e.Emit("GT")
		case ">=":
			e.Emit("GTE")
		default:
			panic("unknown operator: " + exp.Operator)
		}
	case *Identifier:
		idx := e.resolveArgIndex(exp.Name)
		e.Emit("LOAD_ARG", idx)
	case *UInt64Literal:
		e.Emit("PUSH", exp.Value)
	case *StorageAccess:
		e.EmitExpression(exp.Key)
		e.Emit(STORAGE_LOAD, exp.Var)
	default:
		panic(fmt.Sprintf("unsupported expression: %T", exp))
	}
}

func (e *BytecodeEmitter) resolveArgIndex(name string) int {
	for i, a := range e.currentFunc.Args {
		if a.Name == name {
			return i
		}
	}
	panic("arg not found: " + name)
}

func GenerateBytecode(prog *Program) string {
	emitter := NewBytecodeEmitter()
	for _, contract := range prog.Contracts {

		emitter.Emit("PRAGMA", contract.Version+"\n")

		for _, st := range contract.Storages {
			emitter.Emit(STORAGE_DECL, st.Name+"\n")
		}

		for _, f := range contract.Funcs {
			emitter.EmitFunction(f)
		}

	}
	return strings.Join(emitter.Instructions, "\n")
}
