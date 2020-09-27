package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

func main() {
	links, err := parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinsk1: %v\n", err)
		os.Exit(1)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}

func parse(reader io.Reader) ([]string, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	c := n.FirstChild
	if c != nil {
		links = visit(links, c)
	}
	s := n.NextSibling
	if s != nil {
		links = visit(links, s)
	}
	return links
}
