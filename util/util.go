package util

import (
	"bytes"
	"strconv"
    "unicode"
	"unicode/utf8"
	"crypto/rand"
	"runtime"
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
    /* This is ugly, but it is necessary to check the end of the word since
       multiple codes are provided when a word may be a noun, verb, etc. */
    for ; ri>=0; ri-- {
        if !unicode.In(rval[ri-1], unicode.Latin) {
            continue
        } else {
            break
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

/*
This is a utility for use in debugging
*/

func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
