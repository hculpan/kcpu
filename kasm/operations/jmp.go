package operations

import (
	"fmt"
	"strings"

	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerJmpOp(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) (*AssembledOp, []common.AssemblerError) {
	if len(fields) != 2 {
		return nil, []common.AssemblerError{common.NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	var opcode byte = 0x10
	switch strings.ToUpper(fields[0]) {
	case "JEQ":
		opcode = 0x30
	case "JGT":
		opcode = 0x31
	case "JLT":
		opcode = 0x32
	case "JMP":
		opcode = 0x33
	default:
		return nil, []common.AssemblerError{common.NewAssemblerError(fmt.Sprintf("unrecognized JMP op: '%s'", fields[0]), lineNum)}
	}

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
