package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	s := stringsJoin(",", args...)
	fmt.Println(s)
}

func stringsJoin(sep string, vals ...string) string {
	buf := bytes.NewBuffer(nil)
	if len(vals) > 0 {
		for _, s := range vals[:len(vals)-1] {
			buf.WriteString(s)
			buf.WriteString(sep)
		}
		buf.WriteString(vals[len(vals)-1])
	}
	return buf.String()
}
