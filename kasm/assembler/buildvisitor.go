package assembler

import (
	"fmt"
	"strings"

	"github.com/hculpan/kcpu/kasm/common"
	"github.com/hculpan/kcpu/kasm/operations"
)

type BuildVisitor struct {
	AssembledOps    []operations.AssembledOp
	AssemblerErrors []common.AssemblerError
	Symbols         *SymbolsVisitor
}

func NewBuildVisitor(symbols *SymbolsVisitor) BuildVisitor {
	result := BuildVisitor{
		AssembledOps:    make([]operations.AssembledOp, 0),
		AssemblerErrors: make([]common.AssemblerError, 0),
		Symbols:         symbols,
	}

	return result
}

func (b *BuildVisitor) addOp(op *operations.AssembledOp) {
	if op != nil {
		b.AssembledOps = append(b.AssembledOps, *op)
	}
}

func (b *BuildVisitor) addErrors(errors []common.AssemblerError) bool {
	if len(errors) > 0 {
		b.AssemblerErrors = append(b.AssemblerErrors, errors...)
		return true
	}

	return false
}

func (b *BuildVisitor) addOpAndErrors(op *operations.AssembledOp, errors []common.AssemblerError) bool {
	b.addOp(op)
	return b.addErrors(errors)
}

func removeComments(text string) string {
	for {
		pos := strings.Index(text, ";")
		if pos < 0 {
			break
		}
		text = text[:pos]
	}

	return text
}

func (b *BuildVisitor) ProcessLine(lineText string, lineNum int) bool {
	fields := SplitLine(lineText)
	if len(fields) == 0 || fields[0][0] == ':' || fields[0][0] == '.' { // line with label
		return b.addOpAndErrors(operations.NewNoCode(lineText))
	}

	op := strings.ToUpper(fields[0])

	switch op {
	case "LD":
		return b.addOpAndErrors(operations.AssemblerLdOp(fields, lineNum, lineText, b.Symbols))
	case "ST":
		return b.addOpAndErrors(operations.AssemblerStOp(fields, lineNum, lineText, b.Symbols))
	case "STL":
		return b.addOpAndErrors(operations.AssemblerStOp(fields, lineNum, lineText, b.Symbols))
	case "STH":
		return b.addOpAndErrors(operations.AssemblerStOp(fields, lineNum, lineText, b.Symbols))
	case "HALT":
		return b.addOpAndErrors(operations.AssemblerHalt(fields, lineNum, lineText, b.Symbols))
	case "NOOP":
		return b.addOpAndErrors(operations.AssemblerNoOp(fields, lineNum, lineText))
	case "CMP":
		return b.addOpAndErrors(operations.AssemblerCmpOp(fields, lineNum, lineText, b.Symbols))
	case "JMP", "JEQ", "JGT", "JLT":
		return b.addOpAndErrors(operations.AssemblerJmpOp(fields, lineNum, lineText, b.Symbols))
	case "ADD", "SUB":
		return b.addOpAndErrors(operations.AssemblerMathOp(fields, lineNum, lineText, b.Symbols))
	case "OUT":
		return b.addOpAndErrors(operations.AssemblerOutOp(fields, lineNum, lineText, b.Symbols))
	default:
		return b.addErrors([]common.AssemblerError{common.NewAssemblerError("unknown operation '"+op+"'", lineNum)})
	}
}

func (b *BuildVisitor) ToSrings() []string {
	result := []string{}

	addr := 0
	for _, op := range b.AssembledOps {
		if operations.IsNoCodeOp(op) {
			result = append(result, fmt.Sprintf("                      %s", op.OriginalLine))
		} else {
			result = append(result, fmt.Sprintf("%04X%10s        %s", addr, strings.Replace(op.ToString(), ":", "", -1), op.OriginalLine))
			addr += 4
		}
	}

	return result
}

func (b *BuildVisitor) Errors() []common.AssemblerError {
	return b.AssemblerErrors
}
