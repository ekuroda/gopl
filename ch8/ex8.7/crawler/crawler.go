package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/html"
)

var rootDir string = "../download"

type crawler struct {
	root     string
	hostname string
}

type job struct {
	url   string
	depth int
}

type jobResult struct {
	job       *job
	foundJobs []*job
}

func newCrawler(rootURL string) (*crawler, error) {
	root, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	c := &crawler{
		root:     rootURL,
		hostname: root.Hostname(),
	}
	return c, nil
}

func (c *crawler) process() error {
	path := filepath.Join(rootDir, c.hostname)
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	c.breathFirst()
	return nil
}

func (c *crawler) breathFirst() {
	workChan := make(chan *job, 10)
	resultChan := make(chan *jobResult)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for job := range workChan {
				foundJobs := c.crawl(job)
				resultChan <- &jobResult{job: job, foundJobs: foundJobs}
			}
			wg.Done()
		}()
	}

	var jobq []*job
	jobq = append(jobq, &job{
		url:   c.root,
		depth: 0,
	})
	processing := make(map[string]struct{})
	seen := make(map[string]bool)

	go func() {
		for len(jobq) > 0 || len(processing) > 0 {
			if len(jobq) > 0 {
				j := jobq[0]
				select {
				case workChan <- j:
					jobq = jobq[1:]
					processing[j.url] = struct{}{}
				default:
				}
			}
			select {
			case result, ok := <-resultChan:
				if ok {
					delete(processing, result.job.url)
					for _, j := range result.foundJobs {
						if !seen[j.url] {
							seen[j.url] = true
							jobq = append(jobq, j)
						}
					}
				}
			default:
			}
		}
		close(workChan)
		close(resultChan)
	}()

	wg.Wait()
}

func (c *crawler) crawl(j *job) []*job {
	log.Printf("%d: %s", j.depth, j.url)

	u, err := url.Parse(j.url)
	if err != nil {
		log.Printf("failed to parse url %s: %v", j.url, err)
		return nil
	}

	doc, jobs, err := c.openURL(j)
	if err != nil {
		log.Print(err)
		return nil
	}

	requestURI := u.RequestURI()
	filename := url.QueryEscape(requestURI)
	path := filepath.Join(rootDir, c.hostname, filename)

	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			log.Printf("failed to remove file %s: %v", path, err)
			return jobs
		}
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create file %s: %v", path, err)
		return jobs
	}

	defer file.Close()

	err = html.Render(file, doc)
	if err != nil {
		log.Printf("failed to render %s: %v", j.url, err)
	}

	return jobs
}

func (c *crawler) openURL(j *job) (*html.Node, []*job, error) {
	resp, err := http.Get(j.url)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("getting %s: %s", j.url, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading body %s: %v", j.url, err)
	}
	resp.Body.Close()

	reader := bytes.NewReader(b)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %s as HTML: %v", j.url, err)
	}

	jobs := c.forEachNode(j, doc, resp)

	return doc, jobs, nil
}

func (c *crawler) forEachNode(j *job, n *html.Node, resp *http.Response) []*job {
	var jobs []*job
	jobs = append(jobs, c.visitNode(j, n, resp)...)

	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		jobs = append(jobs, c.forEachNode(j, cn, resp)...)
	}

	return jobs
}

func (c *crawler) visitNode(j *job, n *html.Node, resp *http.Response) []*job {
	if n.Type != html.ElementNode {
		return nil
	}

	var jobs []*job

	if n.Data == "a" || n.Data == "area" || n.Data == "base" || n.Data == "link" {
		var attr []html.Attribute
		for _, a := range n.Attr {
			if a.Key != "href" {
				attr = append(attr, a)
				continue
			}
			link, err := url.Parse(a.Val)
			if err != nil {
				log.Print(err)
				attr = append(attr, a)
				continue
			}

			href := a.Val
			if link.Hostname() == c.hostname {
				//log.Printf("\t> href: %s => %s", href, link.RequestURI())
				a.Val = link.RequestURI()
			}

			abs, err := resp.Request.URL.Parse(href)
			if err == nil {
				if abs.Hostname() == c.hostname {
					jobs = append(jobs, &job{
						url:   abs.String(),
						depth: j.depth + 1,
					})
					// if j.depth <= 2 &&
					// 	!strings.HasSuffix(href, ".gz") &&
					// 	!strings.HasSuffix(href, ".zip") &&
					// 	!strings.HasSuffix(href, ".msi") &&
					// 	!strings.HasSuffix(href, ".pkg") {
					// 	jobs = append(jobs, &job{
					// 		url:   abs.String(),
					// 		depth: j.depth + 1,
					// 	})
					// }
				}
			} else {
				log.Print(err)
			}

			attr = append(attr, a)
		}
		n.Attr = attr
	} else if n.Data == "frame" || n.Data == "iframe" || n.Data == "img" || n.Data == "input" || n.Data == "script" {
		var attr []html.Attribute
		for _, a := range n.Attr {
			if a.Key != "src" {
				attr = append(attr, a)
				continue
			}

			abs, err := resp.Request.URL.Parse(a.Val)
			if err != nil {
				log.Print(err)
				attr = append(attr, a)
				continue
			}

			//log.Printf("\t> src: %s => %s", a.Val, abs.String())
			a.Val = abs.String()
			attr = append(attr, a)
		}
		n.Attr = attr
	}

	return jobs
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("url required")
	}

	crawler, err := newCrawler(os.Args[1])
	if err != nil {
		log.Fatalf("invalid url %s: %v", os.Args[1], err)
	}

	if err = crawler.process(); err != nil {
		log.Fatalf("%s", err)
	}
}
