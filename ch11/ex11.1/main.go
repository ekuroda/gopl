package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"unicode"
	"unicode/utf8"
)

type runeSlice []rune

func (r runeSlice) Len() int           { return len(r) }
func (r runeSlice) Less(i, j int) bool { return r[i] < r[j] }
func (r runeSlice) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func main() {
	err := count(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
		os.Exit(1)
	}
}

func count(in io.Reader, out io.Writer) error {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	ir := bufio.NewReader(in)
	for {
		r, n, err := ir.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}

	runes := make([]rune, 0, len(counts))
	for c := range counts {
		runes = append(runes, c)
	}
	sort.Sort(runeSlice(runes))

	fmt.Fprintf(out, "rune\tcount\n")
	for _, c := range runes {
		fmt.Fprintf(out, "%q\t%d\n", c, counts[c])
	}
	fmt.Fprintf(out, "\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Fprintf(out, "%d\t%d\n", i, n)
		}
	}

	if invalid > 0 {
		fmt.Fprintf(out, "\n%d invalid UTF-8 characters\n", invalid)
	}

	return nil
}
