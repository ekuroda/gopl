package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	doc := xml.NewDecoder(os.Stdin)
	var stack []xml.StartElement
	for {
		tok, err := doc.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if elemStrings, ok := containsAll(stack, os.Args[1:]); ok {
				fmt.Printf("%s: %s\n", strings.Join(elemStrings, " "), tok)
			}
		}
	}
}

func containsAll(x []xml.StartElement, y []string) ([]string, bool) {
	elemStrings := make([]string, 0)
	ok := false
	for len(y) <= len(x) {
		if len(y) == 0 {
			ok = true
			break
		}
		if elemString, ok := isMatch(x[0], y[0]); ok {
			elemStrings = append(elemStrings, elemString)
			y = y[1:]
		} else {
			elemStrings = append(elemStrings, x[0].Name.Local)
		}
		x = x[1:]
	}

	if ok {
		for _, e := range x {
			elemStrings = append(elemStrings, e.Name.Local)
		}
	}
	return elemStrings, ok
}

func isMatch(x xml.StartElement, y string) (string, bool) {
	if x.Name.Local == y {
		return x.Name.Local, true
	}

	for _, attr := range x.Attr {
		if attr.Name.Local == "id" {
			if attr.Value == y {
				return fmt.Sprintf("id='%s'", attr.Value), true
			}
		} else if attr.Name.Local == "class" {
			if strings.Contains(attr.Value, y) {
				return fmt.Sprintf("class='%s'", attr.Value), true
			}
		}
	}

	return "", false
}
