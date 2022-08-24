package operations

import "github.com/hculpan/kcpu/kasm/common"

const HALT_OPCODE = 0xFE

func AssemblerHalt(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	var result AssembledOp
	if len(fields) > 1 {
		result = NewAssembledOp(HALT_OPCODE, 0x01, 0xFF, 0xFF, originalLine)
	} else {
		result = NewAssembledOp(HALT_OPCODE, 0xFE, 0xFE, 0xFE, originalLine)

	}
	return &result, nil
}
