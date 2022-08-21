package assembler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	R0 = iota
	R2
	R3
	R4
	R5
	R6
	R7
)

type AssembleFile struct {
	Filename        string
	OutputDirectory string
	OutputFilename  string
}

func NewAssembleFile() *AssembleFile {
	return &AssembleFile{}
}

func removeComments(text string) string {
	for {
		pos := strings.Index(text, "'")
		if pos < 0 {
			break
		}
		text = text[:pos]
	}

	return text
}

func (a *AssembleFile) Assemble() []AssemblerError {
	result := []AssemblerError{}
	output := []AssembledOp{}

	f, err := os.Open(a.Filename)
	if err != nil {
		return []AssemblerError{NewAssemblerError("Unable to open file: "+err.Error(), 0)}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		lineText := scanner.Text()
		op, errs := a.assembleLine(lineText, line)
		if len(errs) > 0 {
			result = append(result, errs...)
		} else if op != nil {
			output = append(output, *op)
		}
	}

	if err := scanner.Err(); err != nil {
		return []AssemblerError{NewAssemblerError("Error reading file: "+err.Error(), 0)}
	}

	a.writeListFile(output)
	a.writeAssembledFile(output)

	return result
}

func (a *AssembleFile) writeListFile(ops []AssembledOp) error {
	listFilename := a.Filename[:len(a.Filename)-len(filepath.Ext(a.Filename))] + ".list"

	f, err := os.Create(listFilename)
	if err != nil {
		return errors.New("unable to create list file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	addr := 0
	for _, op := range ops {
		var err error
		if IsNoCodeOp(op) {
			_, err = w.WriteString(fmt.Sprintf("\t\t\t\t\t\t%s\n", op.OriginalLine))
		} else {
			_, err = w.WriteString(fmt.Sprintf("%04X\t%s\t\t%s\n", addr, strings.Replace(op.ToString(), ":", "", -1), op.OriginalLine))
			addr += 4
		}
		if err != nil {
			return errors.New("unable to write listing file")
		}
	}
	w.Flush()

	return nil
}

func (a *AssembleFile) writeAssembledFile(ops []AssembledOp) error {
	file, err := os.Create(a.OutputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, op := range ops {
		if !IsNoCodeOp(op) {
			bytes := []byte{0, 0, 0, 0}
			bytes[0] = op.Op
			bytes[1] = op.Register
			bytes[2] = op.DataH
			bytes[3] = op.DataL
			_, err := file.Write(bytes)
			if err != nil {
				return errors.New("unable to write to file")
			}
		}
	}

	return nil
}

func Split(r rune) bool {
	return r == ',' || r == ' '
}

func (a *AssembleFile) assembleLine(lineText string, lineNum int) (*AssembledOp, []AssemblerError) {
	text := strings.Trim(removeComments(lineText), " \t\n\r")
	if len(text) == 0 {
		result, _ := NewNoCode(lineText)
		return result, []AssemblerError{}
	}

	fields := strings.FieldsFunc(text, Split)
	if len(fields) == 0 {
		return nil, []AssemblerError{NewAssemblerError("Invalid line", lineNum)}
	}

	op := strings.ToUpper(fields[0])
	switch op {
	case "LD":
		return AssemblerLdOp(fields, lineNum, lineText)
	case "ST":
		return AssemblerStOp(fields, lineNum, lineText)
	case "STL":
		return AssemblerStOp(fields, lineNum, lineText)
	case "STH":
		return AssemblerStOp(fields, lineNum, lineText)
	case "HALT":
		return AssemblerHalt(fields, lineNum, lineText)
	case "NOOP":
		return AssemblerNoOp(fields, lineNum, lineText)
	case "CMP":
		return AssemblerCmpOp(fields, lineNum, lineText)
	case "BRA", "BEQ", "BGT", "BLT":
		return AssemblerBraOp(fields, lineNum, lineText)
	case "ADD", "SUB":
		return AssemblerMathOp(fields, lineNum, lineText)
	case "OUT":
		return AssemblerOutOp(fields, lineNum, lineText)
	default:
		return nil, []AssemblerError{NewAssemblerError("unknown operation '"+op+"'", lineNum)}
	}
}
