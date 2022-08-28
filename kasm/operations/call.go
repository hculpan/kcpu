package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerCallOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 2 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have one argument", lineNum)}
	}

	var opcode byte = 0x70

	var result AssembledOp
	if IsAddress(fields[1]) {
		value, err := FieldToAddress(fields[1])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[1]), lineNum)}
		}
		result = NewAssembledOpWithAddress(opcode, 0x00, uint16(value), originalLine)
	} else if symbolsTable.Exists(fields[1]) {
		value, err := symbolsTable.GetValue(fields[1])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid symbol: '%s'", fields[1]), lineNum)}
		}
		result = NewAssembledOpWithAddress(opcode, 0x00, uint16(value), originalLine)
	} else {
		return nil, []common.AssemblerError{common.NewAssemblerError(fields[0]+" op only accepts address in data parameter", lineNum)}
	}

	return &result, nil
}
