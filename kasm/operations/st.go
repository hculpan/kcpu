package operations

import (
	"fmt"
	"strings"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerStOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 3 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have three arguments", lineNum)}
	}

	var opcode byte = 0x10
	switch strings.ToUpper(fields[0]) {
	case "ST":
		opcode = 0x10
	case "STL":
		opcode = 0x11
	case "STH":
		opcode = 0x12
	default:
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("unrecognized ST op: '%s'", fields[0]), lineNum)}
	}

	r, err := RegisterToNumber(fields[1])
	if err != nil {
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid register identifier: '%s'", fields[1]), lineNum)}
	}

	var result AssembledOp
	if IsRegister(fields[2]) {
		return nil, []common.AssemblerError{common.NewAssemblerError("ST op does not accept register in data parameter", lineNum)}
	} else if IsAddress(fields[2]) {
		value, err := FieldToAddress(fields[2])
		if err != nil {
			return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOpWithAddress(opcode, r, uint16(value), originalLine)
	} else {
		return nil, []common.AssemblerError{common.NewAssemblerError("ST op does not accept value in data parameter", lineNum)}
	}

	return &result, nil
}
