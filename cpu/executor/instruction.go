package executor

import (
	"fmt"
)

type Instruction struct {
	Op       byte
	Register byte
	DataH    byte
	DataL    byte
}

func NewInstruction(op byte, register byte, datah byte, datal byte) Instruction {
	return Instruction{
		Op:       op,
		Register: register,
		DataH:    datah,
		DataL:    datal,
	}
}

func (i Instruction) GetDataAsValue() uint16 {
	var result uint16 = uint16(i.DataH) << 8
	result |= uint16(i.DataL)
	return result
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X", i.Op, i.Register, i.DataH, i.DataL)
}
