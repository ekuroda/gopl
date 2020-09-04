package main

import (
	"fmt"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
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
	for i, course := range topoSort(prereqMap) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
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

func topoSort(m map[string]map[string]struct{}) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(key string)

	visitAll = func(key string) {
		items := m[key]
		for item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(item)
				order = append(order, item)
			}
		}
	}
	for key := range m {
		if !seen[key] {
			seen[key] = true
			visitAll(key)
			order = append(order, key)
		}
	}
	return order
}
