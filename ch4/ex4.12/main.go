package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const workerNum = 10

type comic struct {
	Num        int    `json:"num"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type fetchResult struct {
	num   int
	comic *comic
	err   error
}

type comicIndex map[string]map[int]struct{}

const (
	dataDir   = "./data"
	indexPath = dataDir + "/index.json"
)

const usage = `usage:
  xkcd index
  xkcd search someword`

func main() {
	exitWithUsage := func() {
		fmt.Println(usage)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		exitWithUsage()
	}

	switch os.Args[1] {
	case "index":
		createIndex()
	case "search":
		if len(os.Args) < 3 {
			exitWithUsage()
		}
		search(os.Args[2])
	default:
		exitWithUsage()
	}
}

func createIndex() {
	maxNum, err := fetchComicCount()
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("maxNum=%d\n", maxNum)

	done := make(chan bool)
	fetchResultChan := make(chan *fetchResult)
	comicNumChan := make(chan int, workerNum)

	for i := 0; i < workerNum; i++ {
		go processJob(i, comicNumChan, done, fetchResultChan)
	}

	index := comicIndex{}
	var fetchResults []*fetchResult
	num := 1

	go func() {
		for len(fetchResults) < maxNum {
			fr := <-fetchResultChan
			if fr.err != nil {
				log.Printf("failed to fetch comic: num=%d err=%s\n", fr.num, fr.err)
			} else if fr.comic != nil {
				appendToIndex(index, fr.comic)
				outputData(fr.comic)
			}
			fetchResults = append(fetchResults, fr)
		}
		close(fetchResultChan)
		close(done)
	}()

	for num <= maxNum {
		comicNumChan <- num
		fmt.Printf("write comicNumChan %d\n", num)
		num++
	}
	close(comicNumChan)

	<-done

	saveIndex(index)
}

func processJob(workerID int, comicNumChan <-chan int, done <-chan bool, fetchResultChan chan<- *fetchResult) {
	for {
		select {
		case <-done:
			fmt.Printf("%d: done\n", workerID)
			return
		case num, ok := <-comicNumChan:
			if ok {
				fmt.Printf("%d: num=%d\n", workerID, num)
				comic, err := fetchComic(num)
				fetchResultChan <- &fetchResult{num: num, comic: comic, err: err}
			}
		}
	}
}

func fetchComic(comicNum int) (*comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicNum)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var result comic
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fetchComicCount() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch comic count: %s", resp.Status)
	}

	var comic comic
	if err = json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return 0, err
	}
	return comic.Num, nil
}

func search(query string) {
	index := loadIndex()
	nums, ok := index[query]
	if !ok {
		fmt.Printf("comics not found: query=%s\n", query)
		return
	}

	comics := make([]*comic, 0)
	for num := range nums {
		comic := loadData(num)
		comics = append(comics, comic)
	}

	bytes, err := json.MarshalIndent(comics, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal comics: %s", err)
	}

	fmt.Printf("found %d comics\n%s\n", len(nums), string(bytes))
}

func outputData(comic *comic) {
	bytes, err := json.Marshal(*comic)
	if err != nil {
		log.Fatalf("failed to marshal comic: num=%d, err=%s", comic.Num, err)
	}

	path := fmt.Sprintf("%s/%d.json", dataDir, comic.Num)
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create comic data file: num=%d, err=%s", comic.Num, err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatalf("failed to write comic data file: num=%d, err=%s", comic.Num, err)
	}
}

func loadData(num int) *comic {
	path := fmt.Sprintf("%s/%d.json", dataDir, num)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open comic data file: num=%d, err=%s", num, err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read comic data file: num=%d, err=%s", num, err)
	}

	var comic comic
	err = json.Unmarshal(buf, &comic)
	if err != nil {
		log.Fatalf("failed to unmarshal comic: num=%d, err=%s", num, err)
	}

	return &comic
}

func appendToIndex(index comicIndex, comic *comic) {
	scanner := bufio.NewScanner(strings.NewReader(comic.Transcript))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		w := strings.ToLower(scanner.Text())
		nums, ok := index[w]
		if !ok {
			nums = make(map[int]struct{})
			index[w] = nums
		}
		nums[comic.Num] = struct{}{}
	}
}

func saveIndex(index comicIndex) {
	bytes, err := json.Marshal(index)
	if err != nil {
		log.Fatalf("failed to marshal index: %s", err)
	}

	file, err := os.Create(indexPath)
	if err != nil {
		log.Fatalf("failed to create index file: %s", err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatalf("failed to write index file: %s", err)
	}
}

func loadIndex() comicIndex {
	file, err := os.Open(indexPath)
	if err != nil {
		log.Fatalf("failed to open index file: %s", err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read index file: %s", err)
	}

	var index comicIndex
	err = json.Unmarshal(buf, &index)
	if err != nil {
		log.Fatalf("failed to unmarshal index: %s", err)
	}
	return index
}
