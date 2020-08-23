package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

var pc [256]byte

func main() {
	fmt.Printf("%d", sha256BitDiffCount(os.Args[1], os.Args[2]))
}

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func sha256BitDiffCount(a, b string) (count int) {
	count = bitDiffCount(sha256.Sum256([]byte(a)), sha256.Sum256([]byte(b)))
	return
}

func bitDiffCount(a, b [32]uint8) (count int) {
	for i, v := range a {
		count += popCount(v ^ b[i])
	}
	return
}

func popCount(x uint8) int {
	return int(pc[byte(x)])
}
