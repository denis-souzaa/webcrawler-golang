package website

import (
	"denis-souzaa/web-crawler/db"
	"log"
	"net/http"
	"text/template"
)

type DataLinks struct {
	Links []db.VisitedLink
}

func Run() {
	tmpl, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links, err := db.FindAllLinks()
		if err != nil {

		}

		data := DataLinks{Links: links}

		tmpl.Execute(w, data)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
