package compiler

import (
	"fmt" // Importa fmt para el manejo de errores
	"os"

	"github.com/polarysfoundation/ryot/ast"
	"github.com/polarysfoundation/ryot/codegen"
	"github.com/polarysfoundation/ryot/lexer"
	"github.com/polarysfoundation/ryot/parser"
)

const (
	path = "./artifacts/"
)

// CompiledContract representa el resultado de la compilación de un contrato.
type CompiledContract struct {
	Version  string                // Versión del compilador o del formato de bytecode.
	Bytecode []codegen.Instruction // Las instrucciones de bytecode generadas.
	ABI      codegen.ABI           // La Interfaz Binaria de Aplicación del contrato.
}

// Compile toma el código fuente de un contrato como entrada y lo compila,
// generando bytecode, ABI y archivos de salida.
func Compile(input string) (*CompiledContract, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram() // Asume que ParseProgram ya maneja sus propios errores o los propaga.

	if _, err := os.Stat(path); err == nil {
		os.RemoveAll(path)
	}
	os.MkdirAll(path, os.ModePerm)

	// Verifica si el primer statement es un PragmaStatement y obtiene la versión.
	// Esto asume que el primer statement SIEMPRE será el pragma.
	// Deberías añadir una verificación para asegurarte de que es un PragmaStatement.
	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("prgram is empty")
	}

	pragmaStmt, ok := program.Statements[0].(*ast.PragmaStatement)
	if !ok {
		return nil, fmt.Errorf("expected first statement to be a pragma, got %T", program.Statements[0])
	}
	compilerVersion := pragmaStmt.Value

	// Verifica si hay al menos un ClassStatement (el contrato principal).
	foundContract := false
	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.ClassStatement); ok {
			foundContract = true
			break
		}
	}
	if !foundContract {
		return nil, fmt.Errorf("no se encontró ninguna declaración de contrato")

	}

	compiler := &CompiledContract{
		Version: "1.0.0", // Mantén la versión consistente, o lée la del `codegen`.
	}

	// Si el parser tiene errores, deberías verificarlo aquí.
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("errores de parsing: %s", p.Errors())
	}

	if compilerVersion != compiler.Version {
		return nil, fmt.Errorf("compiler version mismatch: expected %s, got %s", compiler.Version, program.Statements[0].(*ast.PragmaStatement).Value)
	}

	g := codegen.New()

	// La función Generate ahora devuelve un error.
	if err := g.Generate(program); err != nil {
		return nil, fmt.Errorf("error de generación de código: %w", err)
	}

	// Manejo de errores para la escritura de archivos.
	if err := g.WriteABI(path + "abi.json"); err != nil {
		return nil, fmt.Errorf("error al escribir ABI: %w", err)
	}
	if err := g.WriteRYC(path + "bytecode.ryc"); err != nil {
		return nil, fmt.Errorf("error al escribir RYC: %w", err)
	}
	if err := g.WriteRYBC(path + "bytecode.rybc"); err != nil {
		return nil, fmt.Errorf("error al escribir RYBC: %w", err)
	}

	compiler.Bytecode = g.GetInstructions()
	compiler.ABI = g.GetABI()

	return compiler, nil
}
