package assembler

import (
	"fmt"
)

func AssemblerCmpOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []AssemblerError) {
	if len(fields) != 3 {
		return nil, []AssemblerError{NewAssemblerError("invalid op: should have two arguments", lineNum)}
	}

	r, err := RegisterToNumber(fields[1])
	if err != nil {
		return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid register identifier: '%s'", fields[1]), lineNum)}
	}

	var result AssembledOp
	if IsRegister(fields[2]) {
		r2, err := RegisterToNumber(fields[2])
		if err != nil {
			return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid register in data: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOp(0x20, r, 0, r2, originalLine)
	} else {
		value, err := FieldToValue(fields[2])
		if err != nil {
			return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid data value: '%s'", fields[2]), lineNum)}
		}
		result = NewAssembledOpWithAddress(0x21, r, uint16(value), originalLine)
	}

	return &result, nil
}
