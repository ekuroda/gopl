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

	"golang.org/x/net/html"
)

var rootDir string = "../download"

type crawler struct {
	urllist  []string
	hostname string
}

func newCrawler(rootURL string) (*crawler, error) {
	root, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	var urllist []string
	urllist = append(urllist, rootURL)

	c := &crawler{
		urllist:  urllist,
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
	seen := make(map[string]bool)
	for len(c.urllist) > 0 {
		items := c.urllist
		c.urllist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				c.crawl(item)
			}
		}
	}
}

func (c *crawler) crawl(urlString string) {
	fmt.Println(urlString)

	u, err := url.Parse(urlString)
	if err != nil {
		log.Printf("failed to parse url %s: %v", urlString, err)
		return
	}

	doc, err := c.openURL(urlString)
	if err != nil {
		log.Print(err)
		return
	}

	requestURI := u.RequestURI()
	filename := url.QueryEscape(requestURI)
	path := filepath.Join(rootDir, c.hostname, filename)

	if _, err = os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			log.Printf("failed to remove file %s: %v", path, err)
			return
		}
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("failed to create file %s: %v", path, err)
		return
	}

	defer file.Close()

	err = html.Render(file, doc)
	if err != nil {
		log.Printf("failed to render %s: %v", urlString, err)
	}
}

func (c *crawler) openURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body %s: %v", url, err)
	}
	resp.Body.Close()

	reader := bytes.NewReader(b)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	c.forEachNode(doc, resp)

	return doc, nil
}

func (c *crawler) forEachNode(n *html.Node, resp *http.Response) {
	c.visitNode(n, resp)

	for cn := n.FirstChild; cn != nil; cn = cn.NextSibling {
		c.forEachNode(cn, resp)
	}
}

func (c *crawler) visitNode(n *html.Node, resp *http.Response) {
	if n.Type != html.ElementNode {
		return
	}

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
				log.Printf("\t> href: %s => %s", href, link.RequestURI())
				a.Val = link.RequestURI()
			}

			abs, err := resp.Request.URL.Parse(href)
			if err == nil {
				if abs.Hostname() == c.hostname {
					c.urllist = append(c.urllist, abs.String())
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

			log.Printf("\t> src: %s => %s", a.Val, abs.String())
			a.Val = abs.String()
			attr = append(attr, a)
		}
		n.Attr = attr
	}
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
