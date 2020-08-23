package main

import (
	"bufio"
	"fmt"
	"os"
)

type aggregation map[string]*aggregationItem

type aggregationItem struct {
	Count int
	Files fileSet
}

type fileSet map[string]struct{}

func main() {
	agg := make(aggregation)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, agg)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, agg)
			f.Close()
		}
	}

	for line, item := range agg {
		if item.Count > 0 {
			fmt.Printf("%d\t%s\n", item.Count, line)
			for fileName := range item.Files {
				fmt.Printf(":%s\n", fileName)
			}
		}
	}
}

func countLines(f *os.File, aggregation aggregation) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		item, ok := aggregation[line]
		if !ok {
			item = &aggregationItem{Count: 0, Files: fileSet{}}
			aggregation[line] = item
		}
		item.Count++
		item.Files[f.Name()] = struct{}{}
	}
}
