package operations

import (
	"fmt"

	"github.com/hculpan/kcpu/kasm/common"
)

const ORIGIN_NOCODE = 0x01

func AssemblerOriginDirective(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 2 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have one argument", lineNum)}
	}

	var result AssembledOp
	if IsAddress(fields[1]) {
		value, err := FieldToAddress(fields[1])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[1]), lineNum)}
		}
		result = NewAssembledOpWithAddress(NOCODE_OPCODE, ORIGIN_NOCODE, uint16(value), originalLine)
	} else {
		return nil, []common.AssemblerError{common.NewAssemblerError(fields[0]+" op only accepts address in data parameter", lineNum)}
	}

	return &result, nil
}
