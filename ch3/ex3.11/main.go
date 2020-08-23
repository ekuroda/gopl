package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(commaFloatingPointNumberWithSign(os.Args[1]))
}

func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func commaByLoop(s string) string {
	var buf bytes.Buffer
	l := len(s)
	top := 0
	step := l % 3
	if step == 0 {
		step = 3
	}
	buf.WriteString(s[top:step])
	top += step
	step = 3

	for top < l {
		buf.WriteString("," + s[top:top+step])
		top += step
	}
	return buf.String()
}

func commaFloatingPointNumberWithSign(s string) string {
	var buf bytes.Buffer
	integerPartFirstIndex := 0
	if s[0] == '+' || s[0] == '-' {
		buf.WriteString(s[:1])
		integerPartFirstIndex = 1
	}
	dotIndex := strings.Index(s, ".")
	integerPartLastIndex := len(s) - 1
	decimalPartFirstIndex := len(s)
	if dotIndex > 0 {
		integerPartLastIndex = dotIndex - 1
		decimalPartFirstIndex = dotIndex + 1
	}

	integerPart := comma(s[integerPartFirstIndex : integerPartLastIndex+1])
	buf.WriteString(integerPart)
	if dotIndex > 0 {
		buf.WriteString(".")
	}
	if decimalPartFirstIndex < len(s) {
		buf.WriteString(s[decimalPartFirstIndex:len(s)])
	}
	return buf.String()
}
