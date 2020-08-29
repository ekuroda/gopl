package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}

	for _, t := range visit(nil, doc) {
		fmt.Println(t)
	}
}

func visit(texts []string, n *html.Node) []string {
	visitChild := true
	if n.Type == html.TextNode {
		if strings.TrimSpace(n.Data) != "" {
			texts = append(texts, n.Data)
		}
	} else if n.Type == html.ElementNode {
		tag := n.Data
		if tag == "script" || tag == "style" {
			visitChild = false
		}
	}
	if visitChild {
		if c := n.FirstChild; c != nil {
			texts = visit(texts, c)
		}
	}
	if s := n.NextSibling; s != nil {
		texts = visit(texts, s)
	}
	return texts
}
