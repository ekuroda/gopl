package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// Node ...
type Node interface{}

// CharData ...
type CharData string

// Element ...
type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

var d int = 0

func (c CharData) String() string {
	return fmt.Sprintf("%*s%s", d, "", string(c))
}

func (e *Element) String() string {
	attrs := make([]string, 0)
	children := make([]string, 0)

	for _, attr := range e.Attr {
		attrs = append(attrs, fmt.Sprintf("%s: %s", attr.Name.Local, attr.Value))
	}

	d++
	for _, child := range e.Children {
		children = append(children, fmt.Sprintf("%s", child))
	}
	d--

	return fmt.Sprintf("%*sType: %s, Attr: {%s}, Children=[\n%s\n%*s]\n",
		d, "", e.Type.Local, strings.Join(attrs, ", "), strings.Join(children, "\n"), d, "")
}

func main() {
	doc := xml.NewDecoder(os.Stdin)
	var root *Element
	var parentElement, currentElement *Element
	for {
		tok, err := doc.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			elem := &Element{
				Type:     tok.Name,
				Attr:     tok.Attr,
				Children: make([]Node, 0),
			}
			if currentElement == nil {
				currentElement = elem
				root = elem
			} else {
				currentElement.Children = append(currentElement.Children, elem)
				parentElement = currentElement
				currentElement = elem
			}
		case xml.EndElement:
			currentElement = parentElement
		case xml.CharData:
			if currentElement != nil {
				elem := CharData(fmt.Sprintf("%s", tok))
				currentElement.Children = append(currentElement.Children, elem)
			}
		}
	}

	fmt.Printf("%s", root)
}
