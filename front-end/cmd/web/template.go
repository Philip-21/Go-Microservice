package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) renderHome(w http.ResponseWriter, t string) {

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
