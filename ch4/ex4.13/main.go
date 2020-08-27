package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type searchResult struct {
	Search       []*movie `json:"Search"`
	TotalResults string   `json:"totalResults"`
}

type movie struct {
	Title  string `json:"Title"`
	ImdbID string `json:"imdbID"`
	Poster string `json:"Poster"`
}

var (
	apiKey   = os.Getenv("API_KEY")
	isSearch = flag.Bool("s", false, "search movie by word")
	dataURL  = "http://www.omdbapi.com/"
	imageURL = "http://img.omdbapi.com/"
)

const usage = `usage:
  omdb someid
  omdb -s someword`

func main() {
	flag.Parse()

	args := flag.Args()
	if *isSearch {
		if len(args) < 1 {
			exitWithUsage()
		}
		result, err := search(args[0])
		if err != nil {
			log.Fatalf("failed to search: word=%s, error=%s", args[0], err)
		}
		printSearchResult(result)
		return
	}

	if len(args) < 1 {
		exitWithUsage()
	}
	path, err := downloadPoster(args[0])
	if err != nil {
		log.Fatalf("failed to download poster: imdbID=%s, error=%s", args[0], err)
	}
	log.Printf("download succeded. path=%s", path)
}

func exitWithUsage() {
	fmt.Println(usage)
	os.Exit(1)
}

func search(word string) (*searchResult, error) {
	queryWord := url.QueryEscape(word)
	url := fmt.Sprintf("%s?apiKey=%s&s=%s", dataURL, apiKey, queryWord)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %s", resp.Status)
	}

	var result searchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func printSearchResult(searchResult *searchResult) {
	bytes, _ := json.MarshalIndent(searchResult, "", "  ")
	fmt.Println(string(bytes))
}

func downloadPoster(imdbID string) (string, error) {
	url := fmt.Sprintf("%s?apiKey=%s&i=%s", imageURL, apiKey, imdbID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %s", resp.Status)
	}

	path := fmt.Sprintf("%s.jpg", imdbID)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return path, nil
}
