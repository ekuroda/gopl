package main

import (
	"testing"

	"golang.org/x/net/html"
)

func TestElementsByTagName(t *testing.T) {
	doc := &html.Node{
		Type: html.DocumentNode,
	}
	c := &html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: "test1"},
		},
	}
	c.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: "test2"},
		},
	})
	c.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "h1",
	})
	doc.AppendChild(c)

	c = &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}
	c1 := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: "test3"},
		},
	}
	c1.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: "test4"},
		},
	})
	c.AppendChild(c1)
	doc.AppendChild(c)

	nodes := ElementsByTagName(doc, "img")
	if len(nodes) != 2 {
		t.Fatalf("ElementsByTagName(%v) result len=%d; want %d", doc, len(nodes), 2)
	}
	if nodes[0].Data != "img" || len(nodes[0].Attr) < 1 || nodes[0].Attr[0].Val != "test1" {
		val := ""
		if len(nodes[0].Attr) >= 1 {
			val = nodes[0].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) first node: tag, val = %s, %s; want %s, %s", doc, nodes[0].Data, val, "img", "test1")
	}
	if nodes[1].Data != "img" || len(nodes[1].Attr) < 1 || nodes[1].Attr[0].Val != "test4" {
		val := ""
		if len(nodes[1].Attr) >= 1 {
			val = nodes[1].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) second node: tag, val = %s, %s; want %s, %s", doc, nodes[1].Data, val, "img", "test4")
	}

	nodes = ElementsByTagName(doc, "img", "a")
	if len(nodes) != 4 {
		t.Fatalf("ElementsByTagName(%v) result len=%d; want %d", doc, len(nodes), 4)
	}
	if nodes[0].Data != "img" || len(nodes[0].Attr) < 1 || nodes[0].Attr[0].Val != "test1" {
		val := ""
		if len(nodes[0].Attr) >= 1 {
			val = nodes[0].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) first node: tag, val = %s, %s; want %s, %s", doc, nodes[0].Data, val, "img", "test1")
	}

	if nodes[1].Data != "a" || len(nodes[1].Attr) < 1 || nodes[1].Attr[0].Val != "test2" {
		val := ""
		if len(nodes[1].Attr) >= 1 {
			val = nodes[1].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) second node: tag, val = %s, %s; want %s, %s", doc, nodes[1].Data, val, "a", "test2")
	}

	if nodes[2].Data != "a" || len(nodes[2].Attr) < 1 || nodes[2].Attr[0].Val != "test3" {
		val := ""
		if len(nodes[2].Attr) >= 1 {
			val = nodes[2].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) second node: tag, val = %s, %s; want %s, %s", doc, nodes[1].Data, val, "a", "test3")
	}

	if nodes[3].Data != "img" || len(nodes[3].Attr) < 1 || nodes[3].Attr[0].Val != "test4" {
		val := ""
		if len(nodes[3].Attr) >= 1 {
			val = nodes[3].Attr[0].Val
		}
		t.Errorf("ElementsByTagName(%v) second node: tag, val = %s, %s; want %s, %s", doc, nodes[3].Data, val, "img", "test4")
	}
}
