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

	for _, l := range visit(nil, doc) {
		fmt.Println(l)
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		tag := n.Data
		var key string
		if tag == "a" || tag == "link" {
			key = "href"
		} else if tag == "script" || tag == "img" {
			key = "src"
		}
		if key != "" {
			for _, a := range n.Attr {
				if a.Key == key {
					links = append(links, a.Val)
				}
			}
		}
	}
	if c := n.FirstChild; c != nil {
		links = visit(links, c)
	}
	if s := n.NextSibling; s != nil {
		links = visit(links, s)
	}
	return links
}
