package main

import (
	"bufio"
	"gopl/ch4/ex4.14/github"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type repositoryData struct {
	issues        []github.Issue
	milestones    []github.Milestone
	collaborators []github.User
}

const issuesTemplateText = `
<h2>Issues</h2>
<table>
  <tr style='text-align:left'>
    <th>#</th>
    <th>State</th>
    <th>Title</th>
    <th>User</th>
    <th>Milestone</th>
  </tr>
  {{range .}}
  <tr>
    <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
    <td>{{.State}}</td>
    <td>{{.Title}}</td>
    <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
    <td><a href='{{.Milestone.HTMLURL}}'>{{.Milestone.Title}}</a></td>
  </tr>
  {{end}}
</table>
<br>
`
const milestonesTemplateText = `
<h2>Milestones</h2>
<table>
  <tr style='text-align:left'>
    <th>#</th>
    <th>State</th>
    <th>Title</th>
  </tr>
  {{range .}}
  <tr>
    <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
    <td>{{.State}}</td>
    <td>{{.Title}}</td>
  </tr>
  {{end}}
</table>
<br>
`
const collaboratorsTemplateText = `
<h2>Collaborators</h2>
<table>
  <tr style='text-align:left'>
    <th>Id</th>
    <th>Login</th>
  </tr>
  {{range .}}
  <tr>
	<td><a href='{{.HTMLURL}}'>{{.Id}}</a></td>
    <td>{{.Login}}</td>
  </tr>
  {{end}}
</table>
<br>
`

var (
	token                 = os.Getenv("GITHUB_TOKEN")
	issuesTemplate        *template.Template
	milestonesTemplate    *template.Template
	collaboratorsTemplate *template.Template
)

func main() {
	issuesTemplate = template.Must(template.New("issues").Parse(issuesTemplateText))
	milestonesTemplate = template.Must(template.New("milestones").Parse(milestonesTemplateText))
	collaboratorsTemplate = template.Must(template.New("collaborators").Parse(collaboratorsTemplateText))

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL)

	query := r.URL.Query()
	repo := query.Get("repo")
	if repo == "" {
		log.Print("invalid argument: repo required")
		http.Error(w, "invalid argument: repo required", 400)
		return
	}
	rd, err := getRepositoryData(token, repo)
	//log.Printf("repo=%v, rd=%v", repo, rd)
	if err != nil {
		log.Printf("failed to get repository data: %s", err)
		http.Error(w, "failed to get repository data", 500)
		return
	}

	reader := bufio.NewReader(strings.NewReader(`<html lang="ja"><head><meta charset="utf-8"></head><body>`))
	if _, err = reader.WriteTo(w); err != nil {
		log.Printf("failed to write response: %s", err)
		http.Error(w, "failed to write response", 500)
		return
	}

	if err = issuesTemplate.Execute(w, rd.issues); err != nil {
		log.Printf("failed to write issues template: %s", err)
		http.Error(w, "failed to write response", 500)
		return
	}

	if err = milestonesTemplate.Execute(w, rd.milestones); err != nil {
		log.Printf("failed to write milestones template: %s", err)
		http.Error(w, "failed to write response", 500)
		return
	}

	if err = collaboratorsTemplate.Execute(w, rd.collaborators); err != nil {
		log.Printf("failed to write collaborators template: %s", err)
		http.Error(w, "failed to write response", 500)
		return
	}

	reader = bufio.NewReader(strings.NewReader(`</body></html>`))
	if _, err = reader.WriteTo(w); err != nil {
		log.Printf("failed to write response: %s", err)
		http.Error(w, "failed to write response", 500)
	}
}

func getRepositoryData(token, repo string) (*repositoryData, error) {
	issues, err := github.ListIssues(token, repo)
	if err != nil {
		return nil, err
	}
	milestones, err := github.ListMilestones(token, repo)
	if err != nil {
		return nil, err
	}
	collaborators, err := github.ListCollaborators(token, repo)
	if err != nil {
		return nil, err
	}

	rd := &repositoryData{
		issues:        issues,
		milestones:    milestones,
		collaborators: collaborators,
	}

	return rd, nil
}
