package main

import (
	"denis-souzaa/web-crawler/db"
	"denis-souzaa/web-crawler/website"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

var (
	link   string
	action string
)

type VisitedLink struct {
	Website     string    `bson:"website"`
	Link        string    `bson:"link"`
	VisitedDate time.Time `bson:"visited_date"`
}

func init() {
	flag.StringVar(&link, "url", "https://aprendagolang.com.br", "ponto de partida das visitas")
	flag.StringVar(&action, "action", "website", "qual serviço iniciar")
}

func main() {
	flag.Parse()

	switch action {
	case "website":
		website.Run()
	case "webcrawler":
		done := make(chan bool)
		go visitLink(link)

		<-done
	default:
		fmt.Printf("action %s não reconhecida", action)
	}
}

func visitLink(link string) {

	fmt.Printf("visitando: %s\n", link)

	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[error] status diferente de 200: %d\n", resp.StatusCode)
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
			if err != nil || link.Scheme == "" || link.Scheme == "mailto" {
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
			go visitLink(link.String())
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c)
	}
}
