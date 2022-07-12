package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const portNumber = ":80"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.go.html")
	})

	fmt.Println("Starting front end service on port %s\n", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {

	//gets a slice of strings the page requires
	partials := []string{
		"./cmd/web/templates/base.layout.go.html",
		"./cmd/web/templates/header.partial.go.html",
		"./cmd/web/templates/footer.partial.go.html",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
