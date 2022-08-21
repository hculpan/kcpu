package assembler

import (
	"fmt"
	"strings"
)

func AssemblerBraOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []AssemblerError) {
	if len(fields) != 2 {
		return nil, []AssemblerError{NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	var opcode byte = 0x10
	switch strings.ToUpper(fields[0]) {
	case "BEQ":
		opcode = 0x30
	case "BGT":
		opcode = 0x31
	case "BLT":
		opcode = 0x32
	case "BRA":
		opcode = 0x33
	default:
		return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("unrecognized ST op: '%s'", fields[0]), lineNum)}
	}

	var result AssembledOp
	if IsAddress(fields[1]) {
		value, err := FieldToAddress(fields[1])
		if err != nil {
			return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[1]), lineNum)}
		}
		result = NewAssembledOpWithAddress(opcode, 0x00, uint16(value), originalLine)
	} else {
		return nil, []AssemblerError{NewAssemblerError(fields[0]+" op only accepts address in data parameter", lineNum)}
	}

	return &result, nil
}
