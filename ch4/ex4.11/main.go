package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopl/ch4/ex4.11/github"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	token  = os.Getenv("GITHUB_TOKEN")
	editor = os.Getenv("EDITOR")

	command = flag.String("c", "GET", "command (SEARCH/CREATE/GET/UPDATE/CLOSE)")
	owner   = flag.String("o", "", "repository owner")
	repo    = flag.String("r", "", "repository name")
)

type Issue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	//fmt.Println(token)
	flag.Parse()

	switch *command {
	case "SEARCH":
		searchIssues(flag.Args(), false)
	case "CREATE":
		createIssue()
	case "GET":
		getIssue()
	case "UPDATE":
		updateIssue()
	case "CLOSE":
		closeIssue()
	}
	//SearchIssues(os.Args[1:], true)
}

func searchIssues(terms []string, categorize bool) {
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

func createIssue() {
	issue := &github.Issue{}
	editIssue(issue)
	issue, err := github.CreateIssue(token, *owner, *repo, issue)
	if err != nil {
		log.Fatalf("failed to create issue of repo %s/%s: %s", *owner, *repo, err)
	}

	printIssue(issue)
}

func getIssue() {
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("first argument must be a issnue number.\n")
	}

	number := args[0]

	issue, err := github.GetIssue(*owner, *repo, number)
	if err != nil {
		log.Fatalf("failed to get issue %s of repo %s/%s: %s", number, *owner, *repo, err)
	}

	printIssue(issue)
}

func updateIssue() {
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("first argument must be a issnue number.\n")
	}

	number := args[0]

	issue, err := github.GetIssue(*owner, *repo, number)
	if err != nil {
		log.Fatalf("failed to get issue %s of repo %s/%s: %s", number, *owner, *repo, err)
	}

	editIssue(issue)

	err = github.UpdateIssue(token, *owner, *repo, number, issue)
	if err != nil {
		log.Fatalf("failed to update issue %s of repo %s/%s: %s", number, *owner, *repo, err)
	}
	fmt.Println("ok")
}

func closeIssue() {
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("first argument must be a issnue number.\n")
	}

	number := args[0]
	fmt.Println(number)
	err := github.CloseIssue(token, *owner, *repo, number)
	if err != nil {
		log.Fatalf("failed to close issue %s of repo %s/%s: %s", number, *owner, *repo, err)
	}
	fmt.Println("ok")
}

func printIssue(issue *github.Issue) {
	bytes, _ := json.MarshalIndent(issue, "", "  ")
	fmt.Println(string(bytes))
}

func printIssues(issues []*github.Issue) {
	for _, issue := range issues {
		fmt.Printf("#%-5d %9.9s %.55s\n", issue.Number, issue.User.Login, issue.Title)
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

func editIssue(issue *github.Issue) {
	tf, err := ioutil.TempFile("", "issue")
	if err != nil {
		log.Fatal("failed to open temp file: ", err)
	}
	defer func() {
		tf.Close()
		os.Remove(tf.Name())
	}()
	//fmt.Println(tf.Name())

	editContents := Issue{
		Title: issue.Title,
		Body:  issue.Body,
	}
	bytes, err := json.MarshalIndent(&editContents, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal edit contents: %s", err)
	}
	_, err = tf.Write(bytes)
	if err != nil {
		log.Fatalf("failed to write edit contents to temp file: %s", err)
	}
	tf.Sync()

	if editor == "" {
		log.Fatalf("valid editor path must be set to env EDITOR")
	}
	cmd := &exec.Cmd{
		Path:   editor,
		Args:   []string{editor, tf.Name()},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(tf.Name())
	if err = json.Unmarshal(b, &editContents); err != nil {
		log.Fatal(err)
	}

	issue.Title = editContents.Title
	issue.Body = editContents.Body

	//bytes, err = json.MarshalIndent(issue, "", "  ")
	//log.Println(string(bytes))
}
