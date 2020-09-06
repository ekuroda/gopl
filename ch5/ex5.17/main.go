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

	nodes := ElementsByTagName(doc, os.Args[1:]...)
	for _, node := range nodes {
		fmt.Printf("%s, %v\n", node.Data, node.Attr)
	}
}

// ElementsByTagName は doc から name で指定された tag の node を返す
func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	names := make(map[string]struct{})
	for _, n := range name {
		names[n] = struct{}{}
	}

	return forEachNode(doc, names)
}

func forEachNode(n *html.Node, names map[string]struct{}) []*html.Node {
	var nodes []*html.Node
	if n.Type == html.ElementNode {
		if _, ok := names[n.Data]; ok {
			nodes = append(nodes, n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, forEachNode(c, names)...)
	}

	return nodes
}
