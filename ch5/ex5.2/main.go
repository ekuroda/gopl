package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}

	elmCount := make(map[string]int)
	visit(elmCount, doc)
	for key, count := range elmCount {
		fmt.Printf("%-10s %d\n", key, count)
	}
}

func visit(count map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		count[n.Data]++
	}
	c := n.FirstChild
	if c != nil {
		visit(count, c)
	}
	s := n.NextSibling
	if s != nil {
		visit(count, s)
	}
}
