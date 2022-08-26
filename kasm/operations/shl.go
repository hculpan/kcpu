package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerShlOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 3 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	var r byte

	r, err := RegisterToNumber(fields[1])
	if err != nil {
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register identifier: '%s'", fields[1]), lineNum)}
	}

	var result AssembledOp
	value, err := FieldToValue(fields[2], symbolsTable)
	if err != nil {
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid data value: '%s'", fields[2]), lineNum)}
	}
	result = NewAssembledOpWithAddress(0x80, r, uint16(value), originalLine)

	return &result, nil
}
