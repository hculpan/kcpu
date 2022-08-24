package common

import "fmt"

type AssemblerError struct {
	Message string
	Line    int
}

func NewAssemblerError(msg string, lineNum int) AssemblerError {
	return AssemblerError{
		Message: msg,
		Line:    lineNum,
	}
}

func (a AssemblerError) ToString() string {
	return fmt.Sprintf("[%04d] %s", a.Line, a.Message)
}
