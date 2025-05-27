package ryot

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func CompileToBinary(prog *Program) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Header
	buf.Write([]byte("RYOT"))

	// Func count
	funcCount := byte(0)
	for _, c := range prog.Contracts {
		funcCount += byte(len(c.Funcs))
	}
	buf.WriteByte(funcCount)

	// Funcs
	for _, contract := range prog.Contracts {
		for _, fn := range contract.Funcs {
			// FUNC opcode
			buf.WriteByte(OP_FUNC)
			nameLen := byte(len(fn.Name))
			buf.WriteByte(nameLen)
			buf.WriteString(fn.Name)
			buf.WriteByte(byte(len(fn.Args)))

			// Instrucciones
			for _, stmt := range fn.Body {
				switch s := stmt.(type) {
				case *ReturnStatement:
					err := encodeExpr(buf, s.Expr, fn)
					if err != nil {
						return nil, err
					}
					buf.WriteByte(OP_RETURN)
				}
			}
		}
	}

	return buf.Bytes(), nil
}

func encodeExpr(buf *bytes.Buffer, e Expression, fn *FuncDecl) error {
	switch expr := e.(type) {
	case *BinaryExpr:
		if err := encodeExpr(buf, expr.Left, fn); err != nil {
			return err
		}
		if err := encodeExpr(buf, expr.Right, fn); err != nil {
			return err
		}
		if expr.Operator == "+" {
			buf.WriteByte(OP_ADD)
		}
	case *Identifier:
		for i, a := range fn.Args {
			if a.Name == expr.Name {
				buf.WriteByte(OP_LOAD_ARG)
				buf.WriteByte(byte(i))
				return nil
			}
		}
		return fmt.Errorf("unknown identifier: %s", expr.Name)
	case *UInt64Literal:
		buf.WriteByte(OP_PUSH)
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, expr.Value)
		buf.Write(b)
	default:
		return fmt.Errorf("unsupported expression: %T", e)
	}
	return nil
}
