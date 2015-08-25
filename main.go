package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/html"
)

var (
	repoRegexp = regexp.MustCompile(`(github.com|bitbucket.org)/(\w+)/(\w+)`)
)

func process(node *html.Node, w *csv.Writer) error {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "a" {
			for _, a := range child.Attr {
				if a.Key == "href" {
					data := repoRegexp.FindStringSubmatch(a.Val)
					if len(data) == 4 {
						err := w.Write(data)
						return err
					}
					break
				}
			}
		}
		process(child, w)
	}
	return nil
}

func main() {
	resp, err := http.Get("http://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(os.Stdout)
	err = process(doc, w)
	if err != nil {
		log.Fatal(err)
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	w.Flush()
}
