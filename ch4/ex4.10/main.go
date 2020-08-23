package main

import (
	"fmt"
	"gopl/ch4/ex4.10/github"
	"log"
	"os"
	"time"
)

func main() {
	SearchIssues(os.Args[1:], true)
}

func SearchIssues(terms []string, categorize bool) {
	result, err := github.SearchIssues(terms)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	if categorize {
		printCategorizedIssues(result.Items)
	} else {
		printIssues(result.Items)
	}
}

func printIssues(issues []*github.Issue) {
	for _, item := range issues {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}

func printCategorizedIssues(issues []*github.Issue) {
	now := time.Now()
	aMonthAgo := now.AddDate(0, -1, 0)
	aYearAgo := now.AddDate(0, 0, -1)
	categories := []string{"1m", "1y", "1y+"}
	categorized := make(map[string][]*github.Issue)
	for _, category := range categories {
		categorized[category] = make([]*github.Issue, 0)
	}

	for _, item := range issues {
		var key string
		switch {
		case item.CreatedAt.After(aMonthAgo):
			key = "1m"
		case item.CreatedAt.After(aYearAgo):
			key = "1y"
		default:
			key = "1y+"
		}
		categorized[key] = append(categorized[key], item)
	}

	for i, category := range categories {
		fmt.Printf("[%s]\n", category)
		for _, item := range issues {
			fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
		}
		if i < len(categories)-1 {
			fmt.Print("\n")
		}
	}
}
