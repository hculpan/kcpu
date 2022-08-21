package assembler

const NOOP_OPCODE = 0xFF

func AssemblerNoOp(fields []string, lineNum int, originalLine string) (*AssembledOp, []AssemblerError) {
	var result AssembledOp
	if len(fields) > 1 {
		result = NewAssembledOp(NOOP_OPCODE, 0x01, 0xFF, 0xFF, originalLine)
	} else {
		result = NewAssembledOp(NOOP_OPCODE, 0xFF, 0xFF, 0xFF, originalLine)

	}
	return &result, nil
}
