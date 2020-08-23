package main

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"time"
)

var now string = time.Now().Format("20060102150405")

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s.%s.txt", neturl.QueryEscape(url), now))
	if err != nil {
		ch <- fmt.Sprintf("failed to create file: %v", err)
	}

	nbytes, err := io.Copy(file, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
