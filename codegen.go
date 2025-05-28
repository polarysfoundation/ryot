package ryot

import (
	"fmt"
	"strings"
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
		case *ReturnStatement:
			e.EmitExpression(s.Expr)
			e.Emit("RETURN")
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
		for _, fn := range contract.Funcs {
			emitter.EmitFunction(fn)
		}
	}
	return strings.Join(emitter.Instructions, "\n")
}
