package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	n := 10000
	badTime := measure(bad, n)
	goodTime := measure(good, n)
	fmt.Printf("bad=%.2fs, good=%.2fs, diff=%.2fs\n", badTime, goodTime, badTime-goodTime)
}

func measure(f func(), n int) float64 {
	start := time.Now()
	for i := 0; i < n; i++ {
		f()
	}
	elapsed := time.Since(start).Seconds()
	return elapsed
}

func bad() {
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func good() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
