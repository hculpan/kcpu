package executor

import (
	"errors"
	"fmt"
)

const (
	CMP_NOTSET = iota
	CMP_EQUALS
	CMP_GT
	CMP_LT
)

type Cpu struct {
	PC        uint16
	Registers []uint16
	Memory    []byte
	CmpResult int
}

func NewCpu() Cpu {
	return Cpu{
		PC:        0,
		Registers: make([]uint16, 8),
		Memory:    make([]byte, 65536),
	}
}

func (c *Cpu) ToString() string {
	result := fmt.Sprintf("PC:%04X  %02X:%s                CMP:%02X\n", c.PC, c.Memory[c.PC], c.disassembleInstruction(c.Memory[c.PC]), c.CmpResult)
	result += fmt.Sprintf("  R0:%04X  R1:%04X  R2:%04X  R3:%04X  R4:%04X  R5:%04X  R6:%04X  R7:%04X",
		c.Registers[0], c.Registers[1], c.Registers[2], c.Registers[3], c.Registers[4], c.Registers[5], c.Registers[6], c.Registers[7])
	return result
}

func (c *Cpu) SoftReset() {
	c.PC = 0
	for i := 0; i < 8; i++ {
		c.Registers[i] = 0
	}
}

func (c *Cpu) StartExecution() error {
	c.SoftReset()
	for {
		i := c.NextOperation()
		err := c.Execute(i)
		if err != nil {
			return err
		}
	}
}

func (c *Cpu) NextOperation() Instruction {
	result := Instruction{}
	result.Op = c.Memory[c.PC]
	result.Register = c.Memory[c.PC+1]
	result.DataH = c.Memory[c.PC+2]
	result.DataL = c.Memory[c.PC+3]
	return result
}

/*
func outputMemory(addr uint16, cpu *Cpu) {
	fmt.Printf("%04X: %02X\n", addr, cpu.Memory[addr])
}
*/

func (c *Cpu) get16BitValueAtAddress(addr uint16) uint16 {
	var result uint16
	result = uint16(c.Memory[addr]) << 8
	result |= uint16(c.Memory[addr+1])
	return result
}

func (c *Cpu) set16BitValueAtAddress(value uint16, addr uint16) {
	c.Memory[addr] = byte(value >> 8)
	c.Memory[addr+1] = byte(value)
}

func (c *Cpu) disassembleInstruction(op byte) string {
	switch op {
	case 0x00, 0x01, 0x02:
		return "LD"
	case 0x10:
		return "ST"
	case 0x11:
		return "STL"
	case 0x12:
		return "STH"
	case 0x20, 0x21:
		return "CMP"
	case 0x30:
		return "BEQ"
	case 0x31:
		return "BGT"
	case 0x32:
		return "BLT"
	case 0x33:
		return "BRA"
	case 0x40, 0x42, 0x44:
		return "ADD"
	case 0x41, 0x43, 0x45:
		return "SUB"
	case 0x50:
		return "OUT"
	case 0xFE:
		return "HALT"
	case 0xFF:
		return "NOOP"
	default:
		return "UNK"
	}
}

func (c *Cpu) Execute(i Instruction) error {
	switch i.Op {
	case 0x00: // LD r <= value
		c.Registers[i.Register] = i.GetDataAsValue()
		c.PC = c.PC + 4
	case 0x01: // LD r <= r
		c.Registers[i.Register] = c.Registers[i.DataL]
		c.PC = c.PC + 4
	case 0x02: // LD r <= address
		addr := i.GetDataAsValue()
		c.Registers[i.Register] = c.get16BitValueAtAddress(addr)
		c.PC = c.PC + 4
	case 0x10: // ST r, address
		c.set16BitValueAtAddress(c.Registers[i.Register], i.GetDataAsValue())
		c.PC = c.PC + 4
	case 0x11: // STL r, address
		v := byte(c.Registers[i.Register])
		c.Memory[i.GetDataAsValue()] = v
		c.PC = c.PC + 4
	case 0x12: // STH r, address
		v := byte(c.Registers[i.Register] >> 8)
		c.Memory[i.GetDataAsValue()] = v
		c.PC = c.PC + 4
	case 0x20: // CMP r, r
		v1 := c.Registers[i.Register]
		v2 := c.Registers[i.DataL]
		if v1 == v2 {
			c.CmpResult = CMP_EQUALS
		} else if v1 < v2 {
			c.CmpResult = CMP_LT
		} else {
			c.CmpResult = CMP_GT
		}
		c.PC = c.PC + 4
	case 0x21: // CMP r, value
		v1 := c.Registers[i.Register]
		v2 := i.GetDataAsValue()
		if v1 == v2 {
			c.CmpResult = CMP_EQUALS
		} else if v1 < v2 {
			c.CmpResult = CMP_LT
		} else {
			c.CmpResult = CMP_GT
		}
		c.PC = c.PC + 4
	case 0x30: // BEQ addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_EQUALS {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x31: // BGT addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_GT {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x32: // BLT addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_LT {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x33: // BRA addr
		c.PC = i.GetDataAsValue()
	case 0x40: // ADD r, value
		v := i.GetDataAsValue()
		c.Registers[i.Register] += v
		c.PC = c.PC + 4
	case 0x41: // SUB r, value
		v := i.GetDataAsValue()
		c.Registers[i.Register] -= v
		c.PC = c.PC + 4
	case 0x42: // ADD r, r
		c.Registers[i.Register] += c.Registers[i.DataL]
		c.PC = c.PC + 4
	case 0x43: // SUB r, r
		c.Registers[i.Register] -= c.Registers[i.DataL]
		c.PC = c.PC + 4
	case 0x44: // ADD r, addr
		v := c.get16BitValueAtAddress(i.GetDataAsValue())
		c.Registers[i.Register] += v
		c.PC = c.PC + 4
	case 0x45: // SUB r, addr
		v := c.get16BitValueAtAddress(i.GetDataAsValue())
		c.Registers[i.Register] -= v
		c.PC = c.PC + 4
	case 0x50: // OUT r, value
		port := c.Registers[i.Register]
		if port != 0 {
			return fmt.Errorf("unsupport port %02x", port)
		}
		fmt.Print(string(i.DataL))
		c.PC = c.PC + 4
	case 0x51: // OUT r, r
		port := c.Registers[i.Register]
		if port != 0 {
			return fmt.Errorf("unsupport port %02x", port)
		}
		fmt.Print(string(byte(c.Registers[i.DataL])))
		c.PC = c.PC + 4
	case 0xFE:
		if i.Register == 1 {
			fmt.Println(c.ToString())
		}
		return errors.New("halt")
	case 0xFF: // noop
		if i.Register == 1 {
			fmt.Println(c.ToString())
		}
		c.PC = c.PC + 4
	default:
		return fmt.Errorf("unknown operation: %02X", i.Op)
	}
	return nil
}
