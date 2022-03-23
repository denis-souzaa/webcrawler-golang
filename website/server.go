package website

import (
	"log"
	"net/http"
	"text/template"
)

func Run() {
	tmpl, err := template.ParseFiles("website/templates/index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
