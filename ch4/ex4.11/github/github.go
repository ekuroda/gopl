package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL        = "https://api.github.com"
	SearchIssueURL = baseURL + "/search/issues"
	IssuesURL      = baseURL + "/repos/%s/%s/issues"
	IssueURL       = baseURL + "/repos/%s/%s/issues/%s"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

type IssueStatePatch struct {
	State string `json:"state"`
}

type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(SearchIssueURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func CreateIssue(token, owner, repo string, issue *Issue) (*Issue, error) {
	body, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}

	return post(token, owner, repo, body)
}

func GetIssue(owner, repo, number string) (*Issue, error) {
	url := fmt.Sprintf(IssueURL, owner, repo, number)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateIssue(token, owner, repo, number string, issue *Issue) error {
	body, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	return patch(token, owner, repo, number, body)
}

func CloseIssue(token, owner, repo, number string) error {
	statePatch := IssueStatePatch{State: "closed"}
	body, err := json.Marshal(statePatch)
	if err != nil {
		return err
	}

	return patch(token, owner, repo, number, body)
}

func post(token, owner, repo string, body []byte) (*Issue, error) {
	url := fmt.Sprintf(IssuesURL, owner, repo)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func patch(token, owner, repo, number string, body []byte) error {
	url := fmt.Sprintf(IssueURL, owner, repo, number)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}
	return nil
}
