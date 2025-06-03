package codegen

import "fmt"

// Opcode representa un código de operación en el bytecode de Ryot.
type Opcode byte

// Definición de todos los códigos de operación.
const (
	// Operaciones básicas
	OpConst Opcode = iota // 0x00 - Pone una constante en la pila
	OpAdd                 // 0x01 - Suma los dos elementos superiores de la pila
	OpSub                 // 0x02 - Resta los dos elementos superiores de la pila
	OpMul                 // 0x03 - Multiplica los dos elementos superiores de la pila
	OpDiv                 // 0x04 - Divide los dos elementos superiores de la pila
	OpMod                 // 0x05 - Calcula el módulo de los dos elementos superiores de la pila

	// Memoria
	OpStore  // 0x10 - Almacena un valor en la memoria (variable de contrato)
	OpLoad   // 0x11 - Carga un valor de la memoria (variable de contrato)
	OpMStore // 0x12 - Almacena un valor en la memoria transitoria
	OpMLoad  // 0x13 - Carga un valor de la memoria transitoria

	// Control de flujo
	OpJump     // 0x20 - Salto incondicional
	OpJumpI    // 0x21 - Salto condicional
	OpJumpDest // 0x22 - Marca un destino de salto
	OpCall     // 0x23 - Llama a una función
	OpReturn   // 0x24 - Retorna de una función
	OpRevert   // 0x25 - Revierte la ejecución de la transacción

	// Funciones y contratos
	OpContract // 0x30 - Define el inicio de un contrato
	OpFunc     // 0x31 - Define el inicio de una función
	OpEnd      // 0x32 - Marca el final de una sección (contrato, función, enum, struct, storage)

	// Tipos avanzados
	OpArray  // 0x40 - Operación relacionada con arrays
	OpStruct // 0x41 - Operación relacionada con structs
	OpEnum   // 0x42 - Operación relacionada con enums

	// Operaciones específicas de blockchain
	OpAddress // 0x50 - Carga la dirección del contrato o de una entidad
	OpBalance // 0x51 - Carga el balance de una dirección
	OpCaller  // 0x52 - Carga la dirección del llamador
	OpHash    // 0x53 - Calcula el hash de un valor

	// Metadatos
	OpMeta // 0x60 - Operación para metadatos del programa

	// Operaciones de sistema
	OpCreate       // 0x70 - Crea un nuevo contrato
	OpDelete       // 0x71 - Borra una variable de almacenamiento o dato
	OpSelfDestruct // 0x72 - Auto-destruye el contrato

	// Operaciones lógicas
	OpEq  // 0x80 - Igualdad
	OpLt  // 0x81 - Menor que
	OpGt  // 0x82 - Mayor que
	OpAnd // 0x83 - AND lógico
	OpOr  // 0x84 - OR lógico
	OpNot // 0x85 - NOT lógico
	OpNeq // 0x86 - No igual

	// Operaciones de pila
	OpPop  // 0x90 - Elimina el elemento superior de la pila
	OpDup  // 0x91 - Duplica el elemento superior de la pila
	OpSwap // 0x92 - Intercambia los dos elementos superiores de la pila

	OpCheck    // 0xFF - verify to return
	OpErr      // 0xFE - error handling
	OpCheckEnd // 0xFD - verify to end
	OpJumpEnd  // 0xFC - jump to end
	OpLabel    // 0xFB - label for jump

	OpZeroHash // 0xFA - hash de 32 bytes con valor cero (utilizado para inicializar variables o como valor por defecto en estructuras de datos)
	OpZeroAddr // 0xF9 - address con valor cero
)

// Instruction representa una única instrucción de bytecode.
type Instruction struct {
	Opcode Opcode        // El código de operación.
	Args   []interface{} // Argumentos asociados a la instrucción.
	Raw    string        // Representación en cadena de la instrucción para depuración o listado.
}

// String devuelve la representación en cadena del Opcode para una mejor legibilidad.
func (o Opcode) String() string {
	switch o {
	case OpCheck:
		return "CHECK"
	case OpCheckEnd:
		return "CHECK_END"
	case OpJumpEnd:
		return "JUMP_END"
	case OpErr:
		return "ERR"
	case OpNeq:
		return "NEQ"
	case OpLabel:
		return "LABEL"
	case OpConst:
		return "CONST"
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMul:
		return "MUL"
	case OpDiv:
		return "DIV"
	case OpMod:
		return "MOD"
	case OpStore:
		return "STORE"
	case OpLoad:
		return "LOAD"
	case OpMStore:
		return "MSTORE"
	case OpMLoad:
		return "MLOAD"
	case OpJump:
		return "JUMP"
	case OpJumpI:
		return "JUMPI"
	case OpJumpDest:
		return "JUMPDEST"
	case OpCall:
		return "CALL"
	case OpReturn:
		return "RETURN"
	case OpRevert:
		return "REVERT"
	case OpContract:
		return "CONTRACT"
	case OpFunc:
		return "FUNC"
	case OpEnd:
		return "END"
	case OpArray:
		return "ARRAY"
	case OpStruct:
		return "STRUCT"
	case OpEnum:
		return "ENUM"
	case OpAddress:
		return "ADDRESS"
	case OpBalance:
		return "BALANCE"
	case OpCaller:
		return "CALLER"
	case OpHash:
		return "HASH"
	case OpMeta:
		return "META"
	case OpCreate:
		return "CREATE"
	case OpDelete:
		return "DELETE"
	case OpSelfDestruct:
		return "SELFDESTRUCT"
	case OpEq:
		return "EQ"
	case OpLt:
		return "LT"
	case OpGt:
		return "GT"
	case OpAnd:
		return "AND"
	case OpOr:
		return "OR"
	case OpNot:
		return "NOT"
	case OpPop:
		return "POP"
	case OpDup:
		return "DUP"
	case OpSwap:
		return "SWAP"
	default:
		return fmt.Sprintf("UNKNOWN_OPCODE(0x%x)", byte(o))
	}
}
