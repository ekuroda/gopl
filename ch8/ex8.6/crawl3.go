package main

import (
	"flag"
	"fmt"
	"gopl/ch8/ex8.6/links"
	"log"
)

var depth = flag.Uint("depth", 0, "crawl depth")

type link struct {
	url   string
	depth uint
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	flag.Parse()

	worklist := make(chan []*link)
	unseenLinks := make(chan *link)

	links := make([]*link, len(flag.Args()))
	for i, url := range flag.Args() {
		links[i] = &link{url, 0}
	}

	go func() { worklist <- links }()

	for i := 0; i < 20; i++ {
		go func() {
			for unseenLink := range unseenLinks {
				foundLinks := crawl(unseenLink.url)
				nextDepth := unseenLink.depth + 1
				if len(foundLinks) == 0 || nextDepth > *depth {
					continue
				}

				links := make([]*link, len(foundLinks))
				for i, url := range foundLinks {
					links[i] = &link{url: url, depth: nextDepth}
				}
				go func() { worklist <- links }()
			}
		}()
	}

	seen := make(map[string]uint)
	for list := range worklist {
		for _, l := range list {
			if seenDepth, ok := seen[l.url]; !ok || seenDepth > l.depth {
				seen[l.url] = l.depth
				unseenLinks <- l
			}
		}
	}
}
