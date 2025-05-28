package ryot

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func CompileToBinary(prog *Program) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write([]byte("RYOT"))

	funcCount := byte(0)
	for _, c := range prog.Contracts {
		funcCount += byte(len(c.Funcs))
	}
	buf.WriteByte(funcCount)

	for _, contract := range prog.Contracts {
		for _, fn := range contract.Funcs {
			buf.WriteByte(OP_FUNC)
			buf.WriteByte(byte(len(fn.Name)))
			buf.WriteString(fn.Name)
			buf.WriteByte(byte(len(fn.Args)))

			for _, stmt := range fn.Body {
				switch s := stmt.(type) {
				case *ReturnStatement:
					err := encodeExpr(buf, s.Expr, fn)
					if err != nil {
						return nil, err
					}
					buf.WriteByte(OP_RETURN)
				default:
					return nil, fmt.Errorf("unsupported stmt: %T", s)
				}
			}
		}
	}
	return buf.Bytes(), nil
}

func encodeExpr(buf *bytes.Buffer, expr Expression, fn *FuncDecl) error {
	switch e := expr.(type) {
	case *BinaryExpr:
		if err := encodeExpr(buf, e.Left, fn); err != nil {
			return err
		}
		if err := encodeExpr(buf, e.Right, fn); err != nil {
			return err
		}
		switch e.Operator {
		case "+":
			buf.WriteByte(OP_ADD)
		case "-":
			buf.WriteByte(OP_SUB)
		case "*":
			buf.WriteByte(OP_MUL)
		case "/":
			buf.WriteByte(OP_DIV)
		case "%":
			buf.WriteByte(OP_MOD)
		case "==":
			buf.WriteByte(OP_EQ)
		case "!=":
			buf.WriteByte(OP_NEQ)
		case "<":
			buf.WriteByte(OP_LT)
		case "<=":
			buf.WriteByte(OP_LTE)
		case ">":
			buf.WriteByte(OP_GT)
		case ">=":
			buf.WriteByte(OP_GTE)
		default:
			return fmt.Errorf("unsupported operator: %s", e.Operator)
		}
	case *Identifier:
		for i, a := range fn.Args {
			if a.Name == e.Name {
				buf.WriteByte(OP_LOAD_ARG)
				buf.WriteByte(byte(i))
				return nil
			}
		}
		return fmt.Errorf("unknown identifier: %s", e.Name)
	case *UInt64Literal:
		buf.WriteByte(OP_PUSH)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, e.Value)
		buf.Write(b)
		return nil
	default:
		return fmt.Errorf("unsupported expression: %T", e)
	}
	return nil
}
