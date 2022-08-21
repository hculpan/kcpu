package assembler

const HALT_OPCODE = 0xFE

func AssemblerHalt(fields []string, lineNum int, originalLine string) (*AssembledOp, []AssemblerError) {
	var result AssembledOp
	if len(fields) > 1 {
		result = NewAssembledOp(HALT_OPCODE, 0x01, 0xFF, 0xFF, originalLine)
	} else {
		result = NewAssembledOp(HALT_OPCODE, 0xFE, 0xFE, 0xFE, originalLine)

	}
	return &result, nil
}
