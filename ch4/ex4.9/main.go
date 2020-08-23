package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		w := scanner.Text()
		counts[w]++
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "failed to scan: %v\n", scanner.Err())
		os.Exit(1)
	}

	fmt.Println("count\tword")
	for w, n := range counts {
		fmt.Printf("%d\t%q\n", n, w)
	}
}
