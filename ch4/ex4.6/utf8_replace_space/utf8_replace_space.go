package utf8_replace_space

import (
	"unicode"
	"unicode/utf8"
)

func ReplaceSpace(b []byte) []byte {
	isLastRuneSpace := false
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if unicode.IsSpace(r) {
			inc := 0
			if !isLastRuneSpace {
				b[i] = ' '
				inc = 1
			}
			//fmt.Printf("sp: %q %d\n", string(r), inc)
			copy(b[i+inc:], b[i+size:])
			b = b[:len(b)-size+inc]
			i += inc
			isLastRuneSpace = true
		} else {
			//fmt.Printf("ns: %q\n", string(r))
			i += size
			isLastRuneSpace = false
		}
	}
	return b
}
