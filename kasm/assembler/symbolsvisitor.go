package assembler

import (
	"fmt"
	"strings"

	"github.com/hculpan/kcpu/kasm/common"
	"github.com/hculpan/kcpu/kasm/operations"
)

const (
	SYM_TYPE_LABEL = iota
	SYM_TYPE_CONST
)

type Symbol struct {
	Name  string
	Value uint16
	Type  int
}

type SymbolsVisitor struct {
	Symbols     map[string]Symbol
	errors      []common.AssemblerError
	currAddress uint16
}

func NewSymbolsVisitor() SymbolsVisitor {
	return SymbolsVisitor{
		currAddress: 0,
		errors:      make([]common.AssemblerError, 0),
		Symbols:     make(map[string]Symbol),
	}
}

func NewSymbol(name string, value uint16, symbolType int) Symbol {
	return Symbol{
		Name:  name,
		Value: value,
		Type:  symbolType,
	}
}

func (s *SymbolsVisitor) ProcessLine(lineText string, lineNum int) bool {
	fields, err := SplitLine(lineText)
	if err != nil {
		s.errors = append(s.errors, common.NewAssemblerError(err.Error(), lineNum))
		return true
	}

	if len(fields) > 0 && fields[0][0] == ':' {
		label := fields[0][1:]
		s.Symbols[label] = NewSymbol(label, s.currAddress, SYM_TYPE_LABEL)
	} else if len(fields) > 0 && fields[0][0] == '.' {
		switch strings.ToUpper(fields[0]) {
		case ".CONST":
			return s.constDirective(fields, lineText, lineNum)
		case ".DB":
			return s.dbDirective(fields, lineText, lineNum)
		}
	} else if len(fields) > 0 {
		s.currAddress += 4
	}

	return false
}

func (s *SymbolsVisitor) constDirective(fields []string, lineText string, lineNum int) bool {
	if len(fields) != 3 {
		return s.addError(".const must have 2 parameters, NAME and VALUE", lineNum)
	}

	if !common.IsValidVariable(fields[1]) {
		return s.addError(fmt.Sprintf("invalid const name '%s'", fields[1]), lineNum)
	}

	var symbolsTable common.SymbolsTable = s
	v, err := operations.FieldToValue(fields[2], symbolsTable)
	if err != nil {
		return s.addError(fmt.Sprintf("invalid number constant '%s': %s", fields[1], err.Error()), lineNum)
	}

	s.Symbols[fields[1]] = NewSymbol(fields[1], v, SYM_TYPE_CONST)
	return false
}

func (s *SymbolsVisitor) dbDirective(fields []string, lineText string, lineNum int) bool {
	var memLength int = 0
	for i := 1; i < len(fields); i++ {
		str := fields[i]
		if str[0] == '"' {
			str = strings.Replace(str, "\"", "", -1)
			memLength += len(str)
		} else {
			memLength++
		}
	}
	if memLength%4 != 0 {
		memLength = ((memLength / 4) + 1) * 4
	}
	s.currAddress += uint16(memLength)
	return false
}

func (s *SymbolsVisitor) addError(msg string, lineNum int) bool {
	s.errors = append(s.errors, common.NewAssemblerError(msg, lineNum))
	return true
}

func (s *SymbolsVisitor) Errors() []common.AssemblerError {
	return s.errors
}

func (s *SymbolsVisitor) Exists(symbol string) bool {
	if _, result := s.Symbols[symbol]; result {
		return true
	}

	return false
}

func (s *SymbolsVisitor) GetValue(symbol string) (uint16, error) {
	if v, result := s.Symbols[symbol]; result {
		return v.Value, nil
	}

	return 0, fmt.Errorf("call to non-existent symbol '%s'", symbol)
}

func symbolTypeString(symbolType int) string {
	switch symbolType {
	case 0:
		return "label"
	case 1:
		return "const"
	default:
		return "unkn"
	}
}

func (s *SymbolsVisitor) ToStrings() []string {
	result := []string{}
	for s, v := range s.Symbols {
		result = append(result, fmt.Sprintf("  %10s %20s\t%04X", symbolTypeString(v.Type), s, v.Value))
	}
	return result
}
