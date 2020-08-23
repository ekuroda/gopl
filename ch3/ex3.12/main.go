package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("%t\n", isAnagram(os.Args[1], os.Args[2]))
}

func isAnagram(a string, b string) bool {
	if a == b {
		return false
	}

	countMap := make(map[rune]*[2]int)

	for _, c := range a {
		if countMap[c] == nil {
			countMap[c] = &[2]int{}
		}
		countMap[c][0]++
	}

	for _, c := range b {
		if countMap[c] == nil {
			countMap[c] = &[2]int{}
		}
		countMap[c][1]++
	}

	for _, v := range countMap {
		if v[0] != v[1] {
			return false
		}
	}
	return true
}
