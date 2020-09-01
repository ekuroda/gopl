package main

import (
	"fmt"
	"os"
	"regexp"
)

var pattern = regexp.MustCompile(`\$\w+`)

func main() {
	f := func(s string) string {
		return "[" + s + "]"
	}
	s := expand(os.Args[1], f)
	fmt.Println(s)
}

func expand(s string, f func(string) string) string {
	s = pattern.ReplaceAllStringFunc(s, func(m string) string {
		return f(m[1:])
	})
	return s
}
