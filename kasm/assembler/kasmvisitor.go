package assembler

import "github.com/hculpan/kcpu/kasm/common"

type KasmVisitor interface {
	// returns true if errors, otherwise false
	ProcessLine(lineText string, lineNum int) bool
	Errors() []common.AssemblerError
}
