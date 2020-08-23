package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	counts, invalid, err := charcount(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("category\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}

func charcount(in *bufio.Reader) (counts map[string]int, invalid int, err error) {
	counts = make(map[string]int)
	invalid = 0
	for {
		var r rune
		var n int
		r, n, err = in.ReadRune()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			counts, invalid = nil, 0
			return
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		categorize(counts, r)
	}
	return
}

func categorize(count map[string]int, r rune) {
	if unicode.IsControl(r) {
		count["Control"]++
	}
	if unicode.IsDigit(r) {
		count["Digit"]++
	}
	if unicode.IsGraphic(r) {
		count["Graphic"]++
	}
	if unicode.IsLetter(r) {
		count["Letter"]++
	}
	if unicode.IsLower(r) {
		count["Lower"]++
	}
	if unicode.IsMark(r) {
		count["Mark"]++
	}
	if unicode.IsNumber(r) {
		count["Number"]++
	}
	if unicode.IsPrint(r) {
		count["Print"]++
	}
	if unicode.IsPunct(r) {
		count["Punct"]++
	}
	if unicode.IsSpace(r) {
		count["Space"]++
	}
	if unicode.IsSymbol(r) {
		count["Symbol"]++
	}
	if unicode.IsTitle(r) {
		count["Title"]++
	}
	if unicode.IsUpper(r) {
		count["Upper"]++
	}
}
