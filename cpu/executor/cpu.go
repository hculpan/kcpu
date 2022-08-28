package executor

import (
	"errors"
	"fmt"
)

/****************************************
* Registers
*
* R0-R7 : 16-bit general purpose
* PC    : 16-bit program counter
* SP    : 16-bit stack pointer
*
* Flags : 16-bit register
****************************************/

/****************************************
* Memory map
*
* 0x0000 - 0xEFFB : Empty (61,435 bytes)
* 0xEFFC          : JMP 0xF3FC (1 byte)
* 0xF000 - 0xF400 : Video memory (1,000 bytes)
* 0xF400 - 0xF410 : Keyboard buffer (16 bytes)
* 0xF411 - 0xF7FE : Empty (109 bytes)
* 0xF7FF - 0xFFFF : Stack, starts at 0xFFFF
                    going down (2,048 bytes)
****************************************/

const FLAGS = 0xff

const (
	VIDEO_START = 0xF000
	VIDEO_SIZE  = 1000
	VIDEO_END   = VIDEO_START + VIDEO_SIZE
	VIDEO_ROWS  = 25
	VIDEO_COLS  = 40
)

type CpuConfig struct {
	ProgramFilename string
	StartingAddress uint16
}

type Cpu struct {
	Config CpuConfig

	// Computation-related
	PC        uint16
	SP        uint16
	Registers []uint16
	Memory    []byte
	Cycle     int

	// Video-related
	CursorPosInMemory int
	CursorOn          bool
	CursorStatus      bool
}

func NewCpu(config CpuConfig) Cpu {
	result := Cpu{
		Config:    config,
		PC:        0,
		SP:        0xFFFE,
		Registers: make([]uint16, 256),
		Memory:    make([]byte, 65536),
		Cycle:     0,

		// Video-related
		CursorPosInMemory: VIDEO_START,
		CursorOn:          true,
		CursorStatus:      true,
	}

	for i := 0; i < len(result.Memory); i++ {
		result.Memory[i] = 0x00
	}

	result.LoadRom()

	return result
}

func (c *Cpu) ToString() string {
	result := fmt.Sprintf("PC:%04X  %02X:%-4s   %02X  %02X%02X                           SP:%04X [%0s]\n", c.PC, c.Memory[c.PC],
		c.disassembleInstruction(c.Memory[c.PC]), c.Memory[c.PC+1], c.Memory[c.PC+2], c.Memory[c.PC+3], c.SP, c.getStackString())
	result += fmt.Sprintf("R0:%04X  R1:%04X  R2:%04X  R3:%04X  R4:%04X  R5:%04X  R6:%04X  R7:%04X                        Flags:%016b",
		c.Registers[0], c.Registers[1], c.Registers[2], c.Registers[3], c.Registers[4], c.Registers[5], c.Registers[6], c.Registers[7], c.Registers[FLAGS])
	return result
}

func (c *Cpu) getStackString() string {
	result := ""

	var addr int
	for i := 2; i < 10; i++ {
		addr = int(c.SP) + i
		if addr > 0xFFFF {
			break
		}
		result += fmt.Sprintf("%02X:", c.Memory[uint16(addr)])
	}

	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result
}

func (c *Cpu) LoadRom() {
	c.Memory[VIDEO_START-4] = 0x33
	c.Memory[VIDEO_START-3] = 0x00
	c.set16BitValueAtAddress(VIDEO_START-4, VIDEO_START-2)
}

func (c *Cpu) SoftReset() {
	c.PC = 0
	for i := 0; i < 8; i++ {
		c.Registers[i] = 0
	}

	c.SP = 0xFFFE

	c.CursorToPos(0, 0)
	for i := 0; i < VIDEO_SIZE; i++ {
		c.Memory[i+VIDEO_START] = 0
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
	//fmt.Printf("[%10d] %04x\n", c.Cycle, c.PC)
	i := c.NextOperation()
	err := c.Execute(i)
	if err != nil {
		fmt.Println(c.ToString())
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
		c.CursorPosInMemory++
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
		c.CursorPosInMemory = VIDEO_START + ((VIDEO_COLS * row) + col)
	}
}

func (c *Cpu) PosOfCursor() (int, int) {
	result := (c.CursorPosInMemory - VIDEO_START)
	rows := result / VIDEO_COLS
	cols := result % VIDEO_COLS
	return cols, rows
}

func (c *Cpu) GetVideoCharacterLine(line int) []byte {
	startingAddress := VIDEO_START + (line * VIDEO_COLS)
	result := make([]byte, VIDEO_COLS)
	for i := 0; i < VIDEO_COLS; i++ {
		ch := c.Memory[startingAddress+i]
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
	case 0x01, 0x02, 0x03, 0x04, 0x05:
		return "LD"
	case 0x10, 0x13, 0x14:
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
	case 0x50, 0x51:
		return "OUT"
	case 0x60:
		return "PUSH"
	case 0x61:
		return "POP"
	case 0x70:
		return "CALL"
	case 0x71:
		return "RET"
	case 0x80:
		return "SHL"
	case 0x81:
		return "SHR"
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
	case 0x04: // LD r <= [r]
		addr := c.Registers[i.DataL]
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
		c.CompareValues(v1, v2)
		c.PC = c.PC + 4
	case 0x21: // CMP r, value
		v1 := c.Registers[i.Register]
		v2 := i.GetDataAsValue()
		c.CompareValues(v1, v2)
		c.PC = c.PC + 4
	case 0x30: // JEQ addr
		addr := i.GetDataAsValue()
		if c.Registers[FLAGS]&0b0000000000000011 == 3 {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x31: // JGT addr
		addr := i.GetDataAsValue()
		if c.Registers[FLAGS]&0b0000000000000011 == 2 {
			c.PC = addr
		} else {
			c.PC = c.PC + 4
		}
	case 0x32: // JLT addr
		addr := i.GetDataAsValue()
		if c.Registers[FLAGS]&0b0000000000000011 == 1 {
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
	case 0x60: // PUSH r
		c.Push(c.Registers[i.Register])
		c.PC = c.PC + 4
	case 0x61: // POP r
		v, err := c.Pop()
		if err != nil {
			return err
		}
		c.Registers[i.Register] = v
		c.PC = c.PC + 4
	case 0x70: // CALL addr
		c.Push(c.PC)
		c.PC = i.GetDataAsValue()
	case 0x71: // RET
		v, err := c.Pop()
		if err != nil {
			return err
		}
		c.PC = v + 4
	case 0x80: // SHL r, value
		v := c.Registers[i.Register]
		bits := i.GetDataAsValue() % 16
		c.Registers[i.Register] = v << bits
		c.PC = c.PC + 4
	case 0x81: // SHR r, value
		v := c.Registers[i.Register]
		bits := i.GetDataAsValue() % 16
		c.Registers[i.Register] = v >> bits
		c.PC = c.PC + 4
	case 0xFE:
		if i.Register == 1 {
			fmt.Println(c.ToString())
		}
		return errors.New("halt")
	default:
		return fmt.Errorf("unknown operation: %02X", i.Op)
	}

	c.Cycle++
	if c.CursorOn && c.Cycle%200000 == 0 {
		c.CursorStatus = !c.CursorStatus
	}

	return nil
}

func (c *Cpu) Push(v uint16) {
	c.set16BitValueAtAddress(v, c.SP)
	c.SP -= 2
}

func (c *Cpu) Pop() (uint16, error) {
	sp := uint(c.SP)
	sp += 2
	if sp > 0xFFFF {
		return 0, errors.New("stack underflow")
	}
	c.SP = uint16(sp)
	result := c.get16BitValueAtAddress(c.SP)

	return result, nil
}

func (c *Cpu) CompareValues(v1 uint16, v2 uint16) {
	if v1 == v2 {
		c.Registers[FLAGS] |= 1
		c.Registers[FLAGS] |= 1 << 1
	} else if v1 < v2 {
		c.Registers[FLAGS] |= 1
		c.Registers[FLAGS] &= ^(uint16(1) << 1)
	} else {
		c.Registers[FLAGS] &= ^uint16(1)
		c.Registers[FLAGS] |= 1 << 1
	}
}
