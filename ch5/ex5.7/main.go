package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var depth int

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
		os.Exit(1)
	}

	forEachNode(os.Stdout, doc, startElement, endElement)
}

func forEachNode(w io.Writer, n *html.Node, pre, post func(w io.Writer, n *html.Node)) {
	if pre != nil {
		pre(w, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(w, c, pre, post)
	}

	if post != nil {
		post(w, n)
	}
}

func startElement(w io.Writer, n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		depth++
		fmt.Fprintf(w, "%*s<%s", depth*2, "", n.Data)
		attrs := n.Attr
		if len(attrs) > 0 {
			for _, attr := range attrs {
				fmt.Fprintf(w, " %s='%s'", attr.Key, attr.Val)
			}
		}
		if n.FirstChild != nil || n.Data == "script" {
			fmt.Fprint(w, ">\n")
		} else {
			fmt.Fprint(w, "/>\n")
		}
	case html.CommentNode:
		depth++
		fmt.Fprintf(w, "%*s<!-\n", depth*2, "")
		texts := strings.Split(n.Data, "\n")
		for _, text := range texts {
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			fmt.Fprintf(w, "%*s%s\n", depth*2, "", text)
		}
	case html.TextNode:
		depth++
		texts := strings.Split(n.Data, "\n")
		for _, text := range texts {
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			fmt.Fprintf(w, "%*s%s\n", depth*2, "", text)
		}
	}
}

func endElement(w io.Writer, n *html.Node) {
	switch n.Type {
	case html.ElementNode:
		if n.FirstChild != nil || n.Data == "script" {
			fmt.Fprintf(w, "%*s</%s>\n", depth*2, "", n.Data)
		}
		depth--
	case html.CommentNode:
		fmt.Fprintf(w, "%*s-->\n", depth*2, "")
		depth--
	case html.TextNode:
		depth--
	}
}
