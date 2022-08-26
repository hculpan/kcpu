package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerOutOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 3 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	var r byte

	if v, err := symbolsTable.GetValue(fields[1]); err == nil {
		r = byte(v)
	} else {
		v, err := FieldToValue(fields[1], symbolsTable)
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid value for first parameter: '%s'", fields[1]), lineNum)}
		} else {
			r = byte(v)
		}
	}

	var result AssembledOp
	if IsRegister(fields[2]) {
		r2, err := RegisterToNumber(fields[2])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register in data: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOp(0x51, r, 0, r2, originalLine)
	} else {
		value, err := FieldToValue(fields[2], symbolsTable)
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid data value: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOpWithAddress(0x50, r, uint16(value), originalLine)
	}

	return &result, nil
}
