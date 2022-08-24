package assembler

import (
	"strings"
)

func SplitLine(lineText string) []string {
	text := strings.Trim(removeComments(lineText), " \t\n\r")
	if len(text) == 0 {
		return nil
	}

	return strings.FieldsFunc(text, func(r rune) bool {
		return r == ',' || r == ' '
	})
}
