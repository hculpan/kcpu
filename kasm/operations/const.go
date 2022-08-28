package operations

import "github.com/hculpan/kcpu/kasm/common"

/* This is just a special non-op for the purposes of outputting
* empty and comment lines in the lisiting file
*
* It uses the NOOP opcode, but sets register to 0xFF
 */

func AssemblerConstDirective(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	result := NewAssembledOp(NOCODE_OPCODE, 0xFF, 0, 0, originalLine)
	return &result, nil
}
