package operations

import (
	"github.com/hculpan/kcpu/kasm/common"
)

const NOOP_OPCODE = 0x00

func AssemblerNoOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []common.AssemblerError) {
	var result AssembledOp
	if len(fields) > 1 {
		result = NewAssembledOp(NOOP_OPCODE, 0x01, 0xFF, 0xFF, originalLine)
	} else {
		result = NewAssembledOp(NOOP_OPCODE, 0x00, 0xFF, 0xFF, originalLine)

	}
	return &result, nil
}
