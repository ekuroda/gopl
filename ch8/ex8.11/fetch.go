package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"sync"
	"time"
)

var now string = time.Now().Format("20060102150405")

func main() {
	start := time.Now()
	ch := make(chan string)
	done := make(chan struct{})

	var wg sync.WaitGroup
	for _, url := range os.Args[1:] {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := fetch(url, ch, done); err != nil {
				log.Print(err)
			}
		}(url)
	}

	log.Print(<-ch)
	close(done)
	wg.Wait()

	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string, cancel <-chan struct{}) error {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	req.Cancel = cancel
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s.%s.txt", neturl.QueryEscape(url), now))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer file.Close()

	nbytes, err := io.Copy(file, resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("while reading %s: %v", url, err)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
	return nil
}
