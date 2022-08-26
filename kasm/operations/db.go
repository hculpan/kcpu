package operations

import (
	"github.com/hculpan/kcpu/kasm/common"
)

func AssemblerDbDirective(fields []string, lineNum int, originalLine string, symbolsTable common.SymbolsTable) ([]*AssembledOp, []common.AssemblerError) {
	values := []byte{}
	errors := []common.AssemblerError{}

	for i := 1; i < len(fields); i++ {
		if fields[i][0] == '"' {
			str := fields[i]
			values = append(values, []byte(str[1:len(str)-1])...)
		} else {
			v, err := FieldToValue(fields[i], symbolsTable)
			if err != nil {
				errors = append(errors, common.NewAssemblerError(err.Error(), lineNum))
			}
			values = append(values, byte(v))
		}
	}

	// Pad out values with zeroes until multiple of 4
	for {
		if len(values)%4 == 0 {
			break
		}

		values = append(values, 0)
	}

	result := []*AssembledOp{}
	for i := 0; i < len(values); i += 4 {
		op := NewAssembledOpAsData(values[i], values[i+1], values[i+2], values[i+3], originalLine)
		result = append(result, &op)
	}

	return result, errors
}
