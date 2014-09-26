package util

import (
	"bytes"
	"strconv"
	"unicode/utf8"
	"io"
	"math/rand"
	"time"
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

// Randbo creates a stream of non-crypto quality random bytes
type randbo struct {
	rand.Source
}

// New creates a new random reader with a time source.
func New() io.Reader {
	return NewFrom(rand.NewSource(time.Now().UnixNano()))
}

// NewFrom creates a new reader from your own rand.Source
func NewFrom(src rand.Source) io.Reader {
	return &randbo{src}
}

// Read satisfies io.Reader
func (r *randbo) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := int64(r.Int63())
		for i := 0; i < 8; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}
}
