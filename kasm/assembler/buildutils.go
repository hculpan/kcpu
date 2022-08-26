package assembler

import (
	"errors"
	"unicode"
)

func SplitLine(lineText string) ([]string, error) {
	/*
		text := strings.Trim(removeComments(lineText), " \t\n\r")
		if len(text) == 0 {
			return nil
		}

		return strings.FieldsFunc(text, func(r rune) bool {
			return r == ',' || r == ' '
		})
	*/
	var err error = nil
	result := []string{}
	buff := ""
	inQuote := false

	for i := 0; i < len(lineText); i++ {
		c := lineText[i]
		if !inQuote && (unicode.IsSpace(rune(c)) || c == ',') { // split on whitespace and comma
			if len(buff) > 0 {
				result = append(result, buff)
				buff = ""
			}
		} else if !inQuote && c == ';' { // ignore everything after comment ;
			if len(buff) > 0 {
				result = append(result, buff)
				buff = ""
			}
			break
		} else if !inQuote && c == '"' {
			inQuote = true
			buff += string(c)
		} else if inQuote && c == '"' {
			inQuote = false
			if len(buff) > 0 {
				buff += string(c)
				result = append(result, buff)
				buff = ""
			}
		} else {
			buff += string(c)
		}
	}

	if inQuote {
		err = errors.New("unterminated string: missing closing quote")
	}

	if len(buff) > 0 {
		result = append(result, buff)
	}

	return result, err
}
