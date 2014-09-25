package util

import (
	"unicode/utf8"
	"bytes"
	"strconv"
)

func StringToRuneArray(str string) []rune {
	rval := make([]rune, len(str))
	ri := 0
	for i, w := 0, 0; i < len(str); i += w {
		runeValue, width := utf8.DecodeRuneInString(str[i:])
		w = width
		if runeValue > 0 {
			rval[ri] = runeValue
			ri++
		}
	}
	return rval[:ri]
}

func RuneToString(rstr []rune) string {
	var buffer bytes.Buffer
	for _, v := range rstr {
		buffer.WriteString(strconv.QuoteRuneToASCII(v))
	}
	return buffer.String()
}
