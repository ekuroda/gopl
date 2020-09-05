package main

import (
	"fmt"
	"gopl/ch5/ex5.13/links"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var worklist [][2]string
	for _, item := range os.Args[1:] {
		worklist = append(worklist, [2]string{item, ""})
	}
	breathFirst(crawl, worklist)
}

func breathFirst(f func(item, source string) [][2]string, worklist [][2]string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item[0]] {
				seen[item[0]] = true
				worklist = append(worklist, f(item[0], item[1])...)
			}
		}
	}
}

func crawl(urlString, sourceURLString string) [][2]string {
	fmt.Println(urlString)

	responseHandler := func(resp *http.Response, body []byte) {
		storePage := true
		hostname := resp.Request.URL.Hostname()
		if sourceURLString != "" {
			sourceURL, err := url.Parse(sourceURLString)
			if err != nil {
				log.Print(err)
				return
			}
			if hostname != sourceURL.Hostname() {
				storePage = false
			}
		}

		if !storePage {
			return
		}

		path := fmt.Sprintf("%s/%s", hostname, strings.TrimPrefix(resp.Request.URL.Path, "/"))
		if strings.HasSuffix(path, "/") {
			path = path + "index.html"
		}
		dirs := strings.Split(path, "/")
		filePath := filepath.Join(dirs...)
		dirPath, _ := filepath.Split(filePath)

		if dirPath != "" {
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				log.Printf("failed to create directory %s: %s", dirPath, err)
				return
			}
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("failed to create file %s: %s", filePath, err)
			return
		}
		defer file.Close()

		_, err = file.Write(body)
		if err != nil {
			log.Printf("failed to write file %s: %s", filePath, err)
			return
		}
	}

	urllist, err := links.Extract(urlString, responseHandler)
	if err != nil {
		log.Print(err)
	}

	var list [][2]string
	for _, item := range urllist {
		list = append(list, [2]string{item, urlString})
	}
	return list
}
