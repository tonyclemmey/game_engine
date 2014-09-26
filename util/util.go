package util

import (
	"bytes"
	"strconv"
	"unicode/utf8"
	"crypto/rand"
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

func Rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
