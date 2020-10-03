package main

import (
	"bufio"
	"gopl/ch7/ex7.9/sorting"
	"log"
	"net/http"
	"sort"
	"strings"
)

var sorter *sorting.MultiColumnSort

func main() {
	sorter = sorting.NewMultiColumnSort()

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)

	query := r.URL.Query()
	s := query.Get("s")

	switch s {
	case "Title":
		sorter.SelectTitle()
	case "Artist":
		sorter.SelectArtist()
	case "Album":
		sorter.SelectAlbum()
	case "Year":
		sorter.SelectYear()
	case "Length":
		sorter.SelectLength()
	}

	sort.Sort(sorter)

	reader := bufio.NewReader(strings.NewReader(`<html lang="ja"><head><meta charset="utf-8"></head><body>`))
	if _, err := reader.WriteTo(w); err != nil {
		log.Printf("failed to write response: %s", err)
		http.Error(w, "failed to write response", 500)
		return
	}

	if err := sorting.PrintTracks(w); err != nil {
		http.Error(w, "failed to write response", 500)
		return
	}

	reader = bufio.NewReader(strings.NewReader(`</body></html>`))
	if _, err := reader.WriteTo(w); err != nil {
		log.Printf("failed to write response: %s", err)
		http.Error(w, "failed to write response", 500)
	}
}
