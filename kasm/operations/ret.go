package operations

import (
	"github.com/hculpan/kcpu/kasm/common"
)

const RET_OPCODE = 0x71

func AssemblerRetOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []common.AssemblerError) {
	result := NewAssembledOp(RET_OPCODE, 0x00, 0x00, 0x00, originalLine)
	return &result, nil
}
