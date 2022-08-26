package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerLdOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 3 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	r, err := RegisterToNumber(fields[1])
	if err != nil {
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register identifier: '%s'", fields[1]), lineNum)}
	}

	var result AssembledOp
	if fields[2][0] == '[' {
		if fields[2][len(fields[2])-1:] != "]" {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register in data: '%s'", fields[2]), lineNum)}
		}
		r2, err := RegisterToNumber(fields[2][1 : len(fields[2])-1])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register in data: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOp(4, r, 0, r2, originalLine)
	} else if IsRegister(fields[2]) {
		r2, err := RegisterToNumber(fields[2])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register in data: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOp(2, r, 0, r2, originalLine)
	} else if IsAddress(fields[2]) {
		value, err := FieldToAddress(fields[2])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOpWithAddress(3, r, uint16(value), originalLine)
	} else {
		value, err := FieldToValue(fields[2], symbolsTable)
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid data value: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOpWithAddress(1, r, uint16(value), originalLine)
	}

	return &result, nil
}
