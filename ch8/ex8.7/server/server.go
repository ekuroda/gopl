package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var hostname string
var rootDir string = "../download"

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("hostname required")
	}

	hostname = os.Args[1]

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8000", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)

	urlPath := r.URL.Path
	filePath := url.QueryEscape(urlPath)
	filePath = path.Join(rootDir, hostname, filePath)
	f, err := os.Open(filePath)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer f.Close()

	contentType := "text/html"
	if strings.HasSuffix(filePath, ".css") {
		contentType = "text/css"
	}

	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)
}
