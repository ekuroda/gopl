package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	w, i, err := CountWordsAndImages(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("words: %d, images: %d\n", w, i)
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n.Type == html.TextNode {
		scanner := bufio.NewScanner(strings.NewReader(n.Data))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			words++
		}
	} else if n.Type == html.ElementNode {
		if n.Data == "img" {
			images++
		}
	}

	if c := n.FirstChild; c != nil {
		w, i := countWordsAndImages(c)
		words += w
		images += i
	}
	if s := n.NextSibling; s != nil {
		w, i := countWordsAndImages(s)
		words += w
		images += i
	}
	return
}
