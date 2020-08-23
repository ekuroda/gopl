package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	fmt.Println(commaByLoop(os.Args[1]))
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
