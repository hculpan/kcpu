package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerPopOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 2 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have one argument", lineNum)}
	}

	r, err := RegisterToNumber(fields[1])
	if err != nil {
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register identifier: '%s'", fields[1]), lineNum)}
	}

	result := NewAssembledOp(0x61, r, 0x00, 0x00, originalLine)

	return &result, nil
}
