package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

var depth int
var targetID string
var foundNode *html.Node

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}

	node := ElementByID(doc, os.Args[1])
	if node == nil {
		fmt.Println("not found")
		return
	}
	fmt.Println("found")
}

func forEachNode(w io.Writer, n *html.Node, pre, post func(w io.Writer, n *html.Node) bool) bool {
	stop := false
	if pre != nil {
		stop = pre(w, n)
	}

	if !stop {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if stop = forEachNode(w, c, pre, post); stop {
				break
			}
		}
	}

	if !stop {
		if post != nil {
			stop = post(w, n)
		}
	}

	return stop
}

func startElement(w io.Writer, n *html.Node) bool {
	stop := false
	if n.Type == html.ElementNode {
		var id string
		for _, attr := range n.Attr {
			if attr.Key == "id" {
				id = attr.Val
				if id == targetID {
					foundNode = n
					stop = true
					break
				}
			}
		}
		if id == "" {
			fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		} else {
			fmt.Printf("%*s<%s id='%s'>\n", depth*2, "", n.Data, id)
		}
		depth++
	}
	return stop
}

func endElement(w io.Writer, n *html.Node) bool {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
	return false
}

func ElementByID(doc *html.Node, id string) *html.Node {
	targetID = id
	foundNode = nil
	forEachNode(os.Stdout, doc, startElement, endElement)
	return foundNode
}
