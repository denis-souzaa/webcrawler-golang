package main

import (
	"denis-souzaa/web-crawler/db"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

var (
	visited map[string]bool = map[string]bool{}
)

type VisitedLink struct {
	Website     string    `bson:"website"`
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_date"`
}

func main() {
	visitLink("https://aprendagolang.com.br")
}

func visitLink(link string) {

	fmt.Printf("visitando: %s\n", link)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("status diferente de 200: %d", resp.StatusCode))
	}

	node, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	extractLinks(node)
}

func extractLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" {
				continue
			}

			if db.VisitedLink(link.String()) {
				fmt.Printf("link já visitado %s\n", link)
				continue
			}

			visitedLink := VisitedLink{
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			db.Insert("links", visitedLink)
			visitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}
