package main

import "fmt"

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"linear algebra":        {"calculus"},
	"compilers":             {"data structures", "formal languages", "computer organization"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operationg systems"},
	"operationg systems":    {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for key := range prereqs {
		fmt.Printf("[prereqs of %s]\n", key)
		breathFirst(func(item string) []string {
			fmt.Println(item)
			return prereqs[item]
		}, []string{key})
		fmt.Println()
	}
}

func breathFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}
