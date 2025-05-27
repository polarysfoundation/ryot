package ryot

import (
	"fmt"
	"strings"
)

type BytecodeEmitter struct {
	Instructions []string
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
	e.Emit(fmt.Sprintf("FUNC %s %d", f.Name, len(f.Args)))

	for _, stmt := range f.Body {
		switch s := stmt.(type) {
		case *ReturnStatement:
			e.EmitExpression(s.Expr)
			e.Emit("RETURN")
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
	for i, a := range currentFunc.Args {
		if a.Name == name {
			return i
		}
	}
	panic("arg not found: " + name)
}

var currentFunc *FuncDecl

func GenerateBytecode(prog *Program) string {
	emitter := NewBytecodeEmitter()

	for _, contract := range prog.Contracts {
		for _, fn := range contract.Funcs {
			currentFunc = fn
			emitter.EmitFunction(fn)
		}
	}

	return strings.Join(emitter.Instructions, "\n")
}
