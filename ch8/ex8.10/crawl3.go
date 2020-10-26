package main

import (
	"flag"
	"fmt"
	"gopl/ch8/ex8.10/links"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var depth = flag.Uint("depth", 0, "crawl depth")

type link struct {
	url   string
	depth uint
}

func crawl(url string, cancel <-chan struct{}) []string {
	fmt.Println(url)
	list, err := links.Extract(url, cancel)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	flag.Parse()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	worklist := make(chan []*link)
	unseenLinks := make(chan *link)
	cancel := make(chan struct{})

	links := make([]*link, len(flag.Args()))
	for i, url := range flag.Args() {
		links[i] = &link{url, 0}
	}

	go func() { worklist <- links }()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			isDone := false
			for !isDone {
				select {
				case unseenLink, ok := <-unseenLinks:
					if !ok {
						isDone = true
						break
					}

					foundLinks := crawl(unseenLink.url, cancel)
					nextDepth := unseenLink.depth + 1
					if len(foundLinks) == 0 || nextDepth > *depth {
						break
					}

					links := make([]*link, len(foundLinks))
					for i, url := range foundLinks {
						links[i] = &link{url: url, depth: nextDepth}
					}
					go func() { worklist <- links }()
				case <-cancel:
					isDone = true
				}
			}
		}()
	}

	go func() {
		seen := make(map[string]uint)
		isDone := false
		for !isDone {
			select {
			case list, ok := <-worklist:
				if ok {
					for _, l := range list {
						if seenDepth, ok := seen[l.url]; !ok || seenDepth > l.depth {
							seen[l.url] = l.depth
							unseenLinks <- l
						}
					}
				}
			case <-cancel:
				isDone = true
			}
		}
		close(unseenLinks)
	}()

	<-sigs

	fmt.Println("cancel")
	close(cancel)
	wg.Wait()
}
