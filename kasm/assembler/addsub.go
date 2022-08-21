package assembler

import (
	"fmt"
)

func AssemblerMathOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []AssemblerError) {
	if len(fields) != 3 {
		return nil, []AssemblerError{NewAssemblerError("invalid op: should have three arguments", lineNum)}
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
		if fields[0] == "ADD" {
			result = NewAssembledOp(0x42, r, 0, r2, originalLine)
		} else {
			result = NewAssembledOp(0x43, r, 0, r2, originalLine)
		}
	} else if IsAddress(fields[2]) {
		value, err := FieldToAddress(fields[2])
		if err != nil {
			return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid address: '%s'", fields[2]), lineNum)}
		}
		if fields[0] == "ADD" {
			result = NewAssembledOpWithAddress(0x44, r, uint16(value), originalLine)
		} else {
			result = NewAssembledOpWithAddress(0x45, r, uint16(value), originalLine)
		}
	} else {
		value, err := FieldToValue(fields[2])
		if err != nil {
			return nil, []AssemblerError{NewAssemblerError(fmt.Sprintf("invalid data value: '%s'", fields[2]), lineNum)}
		}
		if fields[0] == "ADD" {
			result = NewAssembledOpWithAddress(0x40, r, uint16(value), originalLine)
		} else {
			result = NewAssembledOpWithAddress(0x41, r, uint16(value), originalLine)
		}
	}

	return &result, nil
}
