package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/polarysfoundation/ryot/ast"
)

// Constantes para los números mágicos y la versión del bytecode.
const (
	RyBCMagicNumber  = "\x52\x59\x42\x43" // RYBC
	RyBCVersionMajor = 0x01
	RyBCVersionMinor = 0x00
)

// Generator es el encargado de transformar el AST en instrucciones de bytecode
// y generar la ABI del contrato.
type Generator struct {
	instructions []Instruction // Lista de instrucciones generadas.
	contractName string        // Nombre del contrato actual.
	abi          ABI           // Interfaz Binaria de Aplicación (ABI) del contrato.
	currentFunc  *ABIFunction  // Puntero a la función ABI actual que se está procesando.
}

// ABIType representa un tipo de dato en la ABI.
type ABIType struct {
	Type string `json:"type"`
}

// ABIFunction representa una función en la ABI.
type ABIFunction struct {
	Name       string    `json:"name"`                      // Nombre de la función.
	Inputs     []ABIType `json:"inputs"`                    // Tipos de los parámetros de entrada.
	Outputs    []ABIType `json:"outputs"`                   // Tipos de los valores de retorno.
	Type       string    `json:"type"`                      // Tipo de elemento ABI (e.g., "function", "constructor").
	StateMut   string    `json:"stateMutability,omitempty"` // Mutabilidad del estado (e.g., "pure", "view", "nonpayable", "payable").
	Visibility string    `json:"visibility,omitempty"`      // Visibilidad de la función (e.g., "public", "private").
}

// ABI es una colección de definiciones de funciones ABI.
type ABI []ABIFunction

// New crea e inicializa un nuevo generador de código.
func New() *Generator {
	return &Generator{
		instructions: make([]Instruction, 0),
		abi:          make(ABI, 0),
	}
}

// GetInstructions devuelve la lista de instrucciones de bytecode generadas.
func (g *Generator) GetInstructions() []Instruction {
	return g.instructions
}

// GetABI devuelve la ABI generada para el contrato.
func (g *Generator) GetABI() ABI {
	return g.abi
}

// emit añade una nueva instrucción a la lista de instrucciones generadas.
// Es la forma preferida para añadir instrucciones Opcode al bytecode.
func (g *Generator) emit(op Opcode, args ...interface{}) {
	g.instructions = append(g.instructions, Instruction{
		Opcode: op,
		Args:   args,
		Raw:    g.formatInstruction(op, args...), // Genera el raw string aquí para consistencia
	})
}

// formatInstruction genera la representación en cadena de una instrucción.
// Esto centraliza la lógica para la columna 'Raw' en el RYC.
func (g *Generator) formatInstruction(op Opcode, args ...interface{}) string {
	switch op {
	case OpMeta:
		return fmt.Sprintf("META       %v", args[0])
	case OpContract:
		return fmt.Sprintf("CONTRACT   %v", args[0])
	case OpEnd:
		// El mensaje de OpEnd dependerá del contexto, por lo que podría necesitar ser más dinámico.
		// Para simplificar, asumiremos que Args[0] ya contiene la descripción adecuada.
		if len(args) > 0 {
			return fmt.Sprintf("END_%v", strings.ToUpper(fmt.Sprintf("%v", args[0])))
		}
		return "END" // O un valor por defecto si no se proporciona argumento
	case OpEnum:
		return fmt.Sprintf("ENUM       %v", args[0])
	case OpStruct:
		return fmt.Sprintf("STRUCT     %v", args[0])
	case OpStore:
		return fmt.Sprintf("STORAGE    %v", args[0])
	case OpLoad:
		return fmt.Sprintf("LOAD       %v", args[0])
	case OpConst:
		// Diferenciar el formato de CONST según el tipo.
		if len(args) > 0 {
			switch v := args[0].(type) {
			case uint64:
				return fmt.Sprintf("CONST_U64  %d", v)
			case string:
				return fmt.Sprintf("CONST_STR  \"%s\"", v)
			case bool:
				return fmt.Sprintf("CONST_BOOL %t", v)
			case int: // Para enteros generales, si se usan
				return fmt.Sprintf("CONST_INT  %d", v)
			default:
				return fmt.Sprintf("CONST      %v", v)
			}
		}
		return "CONST"
	case OpAddress:
		return fmt.Sprintf("ADDRESS    %v", args[0])
	case OpHash:
		return fmt.Sprintf("HASH       %v", args[0])
	case OpArray:
		return fmt.Sprintf("ARRAY      [%d elements]", args[0])
	case OpFunc:
		// El formato de OpFunc puede ser más complejo, como ya lo tienes.
		// Asume que el primer arg es el nombre de la función, y el segundo es el retorno si existe.
		sig := fmt.Sprintf("FUNC       %v", args[0])
		if len(args) > 1 && args[1] != nil && args[1].(string) != "" {
			sig += " -> " + args[1].(string)
		}
		return sig
	case OpDelete:
		return fmt.Sprintf("DELETE     %v", args[0])
	case OpCall:
		// Asume que el primer arg es el nombre de la función (o su string representation) y el segundo es el número de args
		return fmt.Sprintf("CALL       %v (%d args)", args[0], args[1])
	case OpReturn:
		return "RETURN" // Return ya no necesita un argumento 'raw' adicional si se emite directamente
	// Añadir más casos según sea necesario para otras opcodes que requieran formato específico en 'Raw'.
	default:
		// Para la mayoría de los opcodes, solo el nombre es suficiente para el 'Raw'
		return fmt.Sprintf("%-10v %v", op.String(), strings.Trim(fmt.Sprint(args...), "[]"))
	}
}

// Generate recorre el árbol de sintaxis abstracta (AST) y genera las instrucciones
// de bytecode correspondientes. Retorna las instrucciones generadas y un error si ocurre.
func (g *Generator) Generate(node ast.Node) error {
	if node == nil {
		return fmt.Errorf("codegen: el nodo AST es nulo")
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			if err := g.Generate(stmt); err != nil {
				return err
			}
		}
	case *ast.PragmaStatement:
		g.emit(OpMeta, n.Value)
	case *ast.ClassStatement:
		g.contractName = n.Name
		g.emit(OpContract, n.Name)
		for _, stmt := range n.Body {
			if err := g.Generate(stmt); err != nil {
				return err
			}
		}
		g.emit(OpEnd, "CONTRACT") // Más específico para el RYC
	case *ast.EnumStatement:
		g.emit(OpEnum, n.Name)
		for _, value := range n.Values {
			g.emit(OpConst, value) // Los valores de enum son constantes
		}
		g.emit(OpEnd, "ENUM")
	case *ast.StructStatement:
		g.emit(OpStruct, n.Name)
		for _, field := range n.Fields {
			// Podrías necesitar un opcode específico para campos de struct o manejarlos como constantes con tipo.
			// Por ahora, se mantiene el uso de OpConst pero podría ser mejor crear OpField o similar.
			g.emit(OpConst, field.Name, field.Type) // Asumiendo que FIELD puede ser representado así.
		}
		g.emit(OpEnd, "STRUCT")
	case *ast.StorageDeclaration:
		g.emit(OpStore, n.Name)
		for _, param := range n.Params {
			g.emit(OpConst, param.Name, param.Type) // Parámetros de almacenamiento.
		}
		// Si `n.Value` representa un valor inicial complejo, necesitaría ser generado recursivamente.
		// Aquí asumimos que `n.Value.Type` es suficiente para el registro.
		g.emit(OpConst, n.Value.Type) // Tipo del valor de almacenamiento.
		g.emit(OpEnd, "STORAGE")

	case *ast.DeleteStatement:
		g.emit(OpDelete, n.Name)
		// OpEnd para DELETE no es típico, las operaciones DELETE suelen ser atómicas.
		// Si "END_DELETE" es para un bloque de instrucciones, esto debe revisarse.
		// Si es solo para marcar el final de la operación, podrías omitirlo o hacer que OpDelete lo implique.
		// Por ahora, lo mantengo si es un requisito de tu formato de RYC.
		g.emit(OpEnd, "DELETE")

	case *ast.StorageAccessStatement:
		g.emit(OpLoad, n.Name)
		for _, param := range n.Params {
			// El segundo argumento de Param es de tipo *ast.Identifier, no es el valor directamente.
			// Necesitas acceder a `param.Value` si quieres el nombre de la variable.
			// `param` es un *ast.Expression, así que g.Generate(param) sería lo correcto si evalúa a un valor.
			g.emit(OpLoad, param.Value)

		}
		g.emit(OpEnd, "LOAD") // Podría ser innecesario dependiendo del significado de END_LOAD.

	case *ast.FuncStatement:
		funcABI := ABIFunction{
			Name:       n.Name,
			Type:       "function",
			Visibility: "public",
		}
		if !n.Public {
			funcABI.Visibility = "private"
		}
		for _, param := range n.Params {
			funcABI.Inputs = append(funcABI.Inputs, ABIType{Type: param.Type})
		}
		if n.ReturnType.Type != "" {
			funcABI.Outputs = append(funcABI.Outputs, ABIType{Type: n.ReturnType.Type})
		}

		g.abi = append(g.abi, funcABI)
		g.currentFunc = &funcABI // Establece la función actual para referencia.

		// Emite la instrucción de función con el nombre y el tipo de retorno.
		g.emit(OpFunc, n.Name, n.ReturnType.Type)

		for _, stmt := range n.Body {
			if err := g.Generate(stmt); err != nil {
				return err
			}
		}

		g.emit(OpEnd, "FUNC")
		g.currentFunc = nil // Limpia la función actual.

	case *ast.ReturnStatement:
		if n.Value != nil {
			if err := g.Generate(n.Value); err != nil {
				return err
			}
		}
		g.emit(OpReturn) // Unificada la emisión de OpReturn.

	case *ast.BinaryExpression:
		if err := g.Generate(n.Left); err != nil {
			return err
		}
		if err := g.Generate(n.Right); err != nil {
			return err
		}
		switch n.Operator {
		case "+":
			g.emit(OpAdd)
		case "-":
			g.emit(OpSub)
		case "*":
			g.emit(OpMul)
		case "/":
			g.emit(OpDiv)
		case "%":
			g.emit(OpMod)
		case "==":
			g.emit(OpEq)
		case "<":
			g.emit(OpLt)
		case ">":
			g.emit(OpGt)
		case "&&":
			g.emit(OpAnd)
		case "||":
			g.emit(OpOr)
		default:
			return fmt.Errorf("codegen: operador binario desconocido '%s'", n.Operator)
		}

	case *ast.IntegerLiteral:
		g.emit(OpConst, uint64(n.Value)) // Usar uint64 para consistencia con RyBC.
	case *ast.StringLiteral:
		g.emit(OpConst, n.Value)
	case *ast.BooleanLiteral:
		g.emit(OpConst, n.Value)
	case *ast.AddressExpression:
		g.emit(OpAddress, n.Value)
	case *ast.HashLiteral:
		g.emit(OpHash, n.Value)
	case *ast.ArrayLiteral:
		// Emite el tamaño del array primero, luego los elementos.
		g.emit(OpArray, len(n.Elements))
		for _, el := range n.Elements {
			if err := g.Generate(el); err != nil {
				return err
			}
		}

	case *ast.CallExpression:
		// Evalúa los argumentos antes de la función.
		for _, arg := range n.Arguments {
			if err := g.Generate(arg); err != nil {
				return err
			}
		}
		// Evalúa la función (que podría ser un identificador o una expresión más compleja).
		if err := g.Generate(n.Function); err != nil {
			return err
		}
		// Pasa la representación de la función y el número de argumentos para el 'Raw'
		g.emit(OpCall, n.Function.String(), len(n.Arguments))

	case *ast.Identifier:
		g.emit(OpLoad, n.Value)

	default:
		return fmt.Errorf("codegen: tipo de nodo AST desconocido para la generación: %T", n)
	}

	return nil // Retorna nil si todo va bien.
}

// WriteABI escribe la ABI generada en un archivo JSON.
func (g *Generator) WriteABI(filename string) error {
	abiData, err := json.MarshalIndent(g.abi, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar ABI: %w", err)
	}
	if err := os.WriteFile(filename, abiData, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo ABI '%s': %w", filename, err)
	}
	return nil
}

// WriteRYC escribe el código de Ryot (bytecode legible por humanos) en un archivo.
func (g *Generator) WriteRYC(filename string, codeHash string) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("; ABI: %s\n", g.contractName))
	builder.WriteString("; Bytecode disassembly\n")
	builder.WriteString("; Source code hash: " + "0x" + codeHash + "\n\n")

	for _, instr := range g.instructions {
		// La propiedad 'Raw' ahora es generada consistentemente por 'emit'.
		builder.WriteString(instr.Raw + "\n")
	}

	if err := os.WriteFile(filename, []byte(builder.String()), 0644); err != nil {
		return fmt.Errorf("error al escribir archivo RYC '%s': %w", filename, err)
	}
	return nil
}

// WriteRYBC escribe el bytecode binario de Ryot en un archivo.
func (g *Generator) WriteRYBC(filename string, codehash []byte) error {
	var bytecode []byte

	// Número mágico para Ryot bytecode (0xRYBC)
	bytecode = append(bytecode, []byte(RyBCMagicNumber)...)

	// Versión (1.0)
	bytecode = append(bytecode, RyBCVersionMajor, RyBCVersionMinor)

	// Añadir el hash del código al bytecode
	if len(codehash) != 32 {
		return fmt.Errorf("codehash debe tener 32 bytes, tiene %d", len(codehash))
	}
	bytecode = append(bytecode, codehash...)

	for _, instr := range g.instructions {
		bytecode = append(bytecode, byte(instr.Opcode))
		// Serializar argumentos
		for _, arg := range instr.Args {
			switch v := arg.(type) {
			case uint64:
				// Usar binary.BigEndian para escribir uint64 de forma segura
				// Necesitarás importar "encoding/binary"
				// O hacerlo manualmente si prefieres mantenerlo sin importaciones adicionales
				/* 	buf := make([]byte, 8) */
				// binary.BigEndian.PutUint64(buf, v) // Esta es la forma profesional
				// bytecode = append(bytecode, buf...)
				// Implementación manual (si no quieres la importación):
				bytecode = append(bytecode, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
					byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			case string:
				bytecode = append(bytecode, []byte(v)...)
				bytecode = append(bytecode, 0x00) // Terminador nulo
			case int:
				// Asumiendo que int es de 4 bytes para este contexto.
				// Podría ser necesario asegurar el tamaño del int o convertir a un tipo de tamaño fijo.
				// buf := make([]byte, 4)
				// binary.BigEndian.PutUint32(buf, uint32(v)) // Si int es 32-bit
				// bytecode = append(bytecode, buf...)
				// Implementación manual:
				bytecode = append(bytecode, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
			case bool:
				if v {
					bytecode = append(bytecode, 0x01)
				} else {
					bytecode = append(bytecode, 0x00)
				}
			// Añadir casos para otros tipos si son posibles argumentos (e.g., float, []byte)
			default:
				// Manejar tipos de argumentos no serializables si es necesario
				return fmt.Errorf("tipo de argumento no serializable en bytecode: %T", v)
			}
		}
	}

	if err := os.WriteFile(filename, bytecode, 0644); err != nil {
		return fmt.Errorf("error al escribir archivo RYBC '%s': %w", filename, err)
	}
	return nil
}
