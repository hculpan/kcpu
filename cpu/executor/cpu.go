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

const (
	VIDEO_START = 62464
	VIDEO_SIZE  = 1920
	VIDEO_END   = VIDEO_START + VIDEO_SIZE
	VIDEO_ROWS  = 24
	VIDEO_COLS  = 40
)

type Cpu struct {
	PC                uint16
	SP                uint16
	Registers         []uint16
	Memory            []byte
	CmpResult         int
	CursorPosInMemory int
	CursorOn          bool
	CursorStatus      bool
	Cycle             int
}

func NewCpu() Cpu {
	result := Cpu{
		PC:                0,
		SP:                0xFFFF,
		Registers:         make([]uint16, 8),
		Memory:            make([]byte, 65536),
		CursorPosInMemory: VIDEO_START,
		CursorOn:          true,
		CursorStatus:      true,
		Cycle:             0,
	}
	result.LoadRom()

	return result
}

func (c *Cpu) ToString() string {
	result := fmt.Sprintf("PC:%04X  %02X:%s                CMP:%02X\n", c.PC, c.Memory[c.PC], c.disassembleInstruction(c.Memory[c.PC]), c.CmpResult)
	result += fmt.Sprintf("  R0:%04X  R1:%04X  R2:%04X  R3:%04X  R4:%04X  R5:%04X  R6:%04X  R7:%04X",
		c.Registers[0], c.Registers[1], c.Registers[2], c.Registers[3], c.Registers[4], c.Registers[5], c.Registers[6], c.Registers[7])
	return result
}

func (c *Cpu) LoadRom() {
	c.Memory[VIDEO_START-4] = 0x33
	c.Memory[VIDEO_START-3] = 0x00
	c.Memory[VIDEO_START-2] = 0xF3
	c.Memory[VIDEO_START-1] = 0xFC
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

func (c *Cpu) ExecuteSingle() error {
	i := c.NextOperation()
	err := c.Execute(i)
	if err != nil {
		return err
	}

	return nil
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

func (c *Cpu) SetVideoToCharacter(col int, row int, ch byte) {
	c.Memory[VIDEO_START+(((row*VIDEO_COLS)+col)*2)] = ch
}

func (c *Cpu) SetCharacterAtCursor(ch byte) {
	switch ch {
	case 10: // newline
		_, rows := c.PosOfCursor()
		rows++
		if rows >= VIDEO_ROWS {
			rows--
		}
		c.CursorToPos(0, rows)
	default:
		c.Memory[c.CursorPosInMemory] = ch
		c.CursorPosInMemory += 2
		if c.CursorPosInMemory >= VIDEO_END {
			c.CursorHome()
		}
	}
}

func (c *Cpu) CursorHome() {
	c.CursorPosInMemory = VIDEO_START
}

func (c *Cpu) CursorToPos(col int, row int) {
	if row >= 0 && row < VIDEO_ROWS && col >= 0 && col < VIDEO_COLS {
		c.CursorPosInMemory = VIDEO_START + (((VIDEO_COLS * row) + col) * 2)
	}
}

func (c *Cpu) PosOfCursor() (int, int) {
	result := (c.CursorPosInMemory - VIDEO_START) / 2
	rows := result / VIDEO_COLS
	cols := result % VIDEO_COLS
	return cols, rows
}

func (c *Cpu) GetVideoCharacterLine(line int) []byte {
	startingAddress := VIDEO_START + (line * VIDEO_COLS * 2)
	result := make([]byte, 40)
	for i := 0; i < VIDEO_COLS; i++ {
		ch := c.Memory[startingAddress+(i*2)]
		if c.CursorOn && c.CursorStatus {
			ccols, crows := c.PosOfCursor()
			if crows == line && ccols == i {
				ch = '~'
			}
		}
		if ch != 10 && (ch < 32 || ch > 126) {
			ch = 32
		}
		result[i] = ch
	}
	return result
}

func (c *Cpu) disassembleInstruction(op byte) string {
	switch op {
	case 0x00:
		return "NOOP"
	case 0x01, 0x02, 0x03:
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
		return "JEQ"
	case 0x31:
		return "JGT"
	case 0x32:
		return "JLT"
	case 0x33:
		return "JMP"
	case 0x40, 0x42, 0x44:
		return "ADD"
	case 0x41, 0x43, 0x45:
		return "SUB"
	case 0x50:
		return "OUT"
	case 0xFE:
		return "HALT"
	default:
		return "UNK"
	}
}

func (c *Cpu) Execute(i Instruction) error {
	switch i.Op {
	case 0x00: // noop
		if i.Register == 1 {
			fmt.Println(c.ToString())
		}
		c.PC = c.PC + 4
	case 0x01: // LD r <= value
		c.Registers[i.Register] = i.GetDataAsValue()
		c.PC = c.PC + 4
	case 0x02: // LD r <= r
		c.Registers[i.Register] = c.Registers[i.DataL]
		c.PC = c.PC + 4
	case 0x03: // LD r <= address
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
	case 0x30: // JEQ addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_EQUALS {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x31: // JGT addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_GT {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x32: // JLT addr
		addr := i.GetDataAsValue()
		if c.CmpResult == CMP_LT {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x33: // JMP addr
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
		c.SetCharacterAtCursor(i.DataL)
		c.PC = c.PC + 4
	case 0x51: // OUT value, r
		port := i.Register
		if port != 0 {
			return fmt.Errorf("unsupport port %02x", port)
		}
		c.SetCharacterAtCursor(byte(c.Registers[i.DataL]))
		c.PC = c.PC + 4
	case 0xFE:
		if i.Register == 1 {
			fmt.Println(c.ToString())
		}
		return errors.New("halt")
	case 0xFF: // noop
	default:
		return fmt.Errorf("unknown operation: %02X", i.Op)
	}

	c.Cycle++
	if c.CursorOn && c.Cycle%200000 == 0 {
		c.CursorStatus = !c.CursorStatus
	}

	return nil
}
