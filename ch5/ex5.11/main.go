package main

import (
	"fmt"
	"strings"
)

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
	prereqMap := createPrereqMap()
	order, cyclic := topoSort(prereqMap)
	for i, course := range order {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
	fmt.Printf("cyclic items: %s\n", strings.Join(cyclic, ", "))
}

func createPrereqMap() map[string]map[string]struct{} {
	prereqMap := make(map[string]map[string]struct{})
	for key, values := range prereqs {
		m := make(map[string]struct{})
		for _, value := range values {
			m[value] = struct{}{}
		}
		prereqMap[key] = m
	}
	return prereqMap
}

func topoSort(m map[string]map[string]struct{}) ([]string, []string) {
	var order []string
	var cyclic []string
	seen := make(map[string]bool)
	var visitAll func(key string)

	visitItem := func(item string) {
		s, ok := seen[item]
		if s && ok {
			cyclic = append(cyclic, item)
		}
		if !ok {
			seen[item] = true
			visitAll(item)
			seen[item] = false
			order = append(order, item)
		}
	}

	visitAll = func(key string) {
		items := m[key]
		for item := range items {
			visitItem(item)
		}
	}

	for key := range m {
		visitItem(key)
	}
	return order, cyclic
}
