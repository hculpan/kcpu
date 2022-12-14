package operations

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/hculpan/kcpu/kasm/common"
)

type AssembledOp struct {
	Op           byte
	Register     byte
	DataH        byte
	DataL        byte
	OriginalLine string
}

func NewAssembledOp(op byte, register byte, datah byte, datal byte, originalLine string) AssembledOp {
	return AssembledOp{
		Op:           op,
		Register:     register,
		DataH:        datah,
		DataL:        datal,
		OriginalLine: originalLine,
	}
}

func NewAssembledOpWithAddress(op byte, register byte, addr uint16, originalLine string) AssembledOp {
	return AssembledOp{
		Op:           op,
		Register:     register,
		DataH:        byte(addr >> 8),
		DataL:        byte(addr & 0x00FF),
		OriginalLine: originalLine,
	}
}

func NewAssembledOpAsData(data1 byte, data2 byte, data3 byte, data4 byte, originalLine string) AssembledOp {
	return AssembledOp{
		Op:           data1,
		Register:     data2,
		DataH:        data3,
		DataL:        data4,
		OriginalLine: originalLine,
	}
}

func (a AssembledOp) GetDataAsAddress() uint16 {
	var result uint16 = uint16(a.DataH) << 8
	result |= uint16(a.DataL)
	return result
}

func (a AssembledOp) ToString() string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X", a.Op, a.Register, a.DataH, a.DataL)
}

func RegisterToNumber(r string) (byte, error) {
	if r[0] != 'R' && r[0] != 'r' {
		return 0, errors.New("invalid register identifier: must begin with 'R' or 'r'")
	}

	var regNum byte
	if strings.ToUpper(r[1:]) == "FLAGS" {
		regNum = 0xFF
	} else {
		n, err := strconv.Atoi(r[1:])
		if err != nil {
			return 0, fmt.Errorf("invalid register identifer: %s", err)
		} else if n < 0 || n > 255 {
			return 0, errors.New("invalid register identifier: second character must be digit of 0-255")
		} else {
			regNum = byte(n)
		}
	}

	return regNum, nil
}

func IsRegister(r string) bool {
	result, _ := regexp.MatchString("^[Rr]\\d{1,3}$", r)
	return result
}

func IsAddress(r string) bool {
	return r[0] == '$'
}

func IsSymbol(r string, symbolsTable common.SymbolsTable) bool {
	return symbolsTable.Exists(r)
}

func FieldToAddress(r string) (uint16, error) {
	return FieldToValue(r[1:], nil)
}

func FieldToValue(r string, symbolsTable common.SymbolsTable) (uint16, error) {
	if r[0] == '\'' {
		if r[len(r)-1] != '\'' {
			return 0, errors.New("unterminted character")
		} else if len(r) > 3 {
			return 0, errors.New("only one character permitted")
		}

		if len(r) == 2 {
			return 0, nil
		} else {
			return uint16(r[1]), nil
		}
	} else if unicode.IsDigit(rune(r[0])) {
		base := 10
		if len(r) > 2 && (strings.HasPrefix(r, "0x") || strings.HasPrefix(r, "0X")) {
			base = 16
			r = r[2:]
		} else if len(r) > 2 && (strings.HasPrefix(r, "0B") || strings.HasPrefix(r, "0b")) {
			base = 2
			r = r[2:]
		} else if len(r) > 2 && (strings.HasPrefix(r, "0O") || strings.HasPrefix(r, "0o")) {
			base = 8
			r = r[2:]
		}
		result, err := strconv.ParseInt(r, base, 64)
		if err == nil && (result < 0 || result > 65535) {
			err = errors.New("number out of range")
		}
		return uint16(result), err
	} else {
		return symbolsTable.GetValue(r)
	}
}
