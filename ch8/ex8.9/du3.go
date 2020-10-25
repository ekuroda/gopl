package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

type dirData struct {
	name           string
	nfiles, nbytes int64
}

type fileSizeData struct {
	root string
	size int64
}

func main() {
	flag.Parse()

	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	dirDataMap := make(map[string]*dirData)
	dirDataList := make([]*dirData, 0)

	fileSizes := make(chan *fileSizeData)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		dirDataMap[root] = &dirData{name: root}
		dirDataList = append(dirDataList, dirDataMap[root])
		go walkDir(root, root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}

loop:
	for {
		select {
		case data, ok := <-fileSizes:
			if !ok {
				break loop
			}
			fsd, ok := dirDataMap[data.root]
			if ok {
				fsd.nfiles++
				fsd.nbytes += data.size
			}
		case <-tick:
			printDiskUsage(dirDataList)
		}
	}

	printDiskUsage(dirDataList)
}

func printDiskUsage(dirDataList []*dirData) {
	var lines []string
	for _, dirData := range dirDataList {
		lines = append(
			lines,
			fmt.Sprintf("%s: %d files  %.1f GB", dirData.name, dirData.nfiles, float64(dirData.nbytes)/1e9))
	}
	fmt.Println(strings.Join(lines, "\n"))
	fmt.Println()
}

func walkDir(dir, root string, n *sync.WaitGroup, fileSizes chan<- *fileSizeData) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, root, n, fileSizes)
		} else {
			fileSizes <- &fileSizeData{root, entry.Size()}
		}
	}
}

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
