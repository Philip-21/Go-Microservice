package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignUp struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
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
	log.Println("main page gotten")
}

func (app *Config) renderAuth(w http.ResponseWriter, v string) error {
	render := "./cmd/web/templates/auth.page.go.html"
	t, err := template.ParseFiles(render)
	if err != nil {
		log.Println(err)
		return err
	}
	data := AuthPayload{
		Email:    "",
		Password: "",
	}
	err = t.ExecuteTemplate(w, "auth", data)
	if err != nil {
		log.Printf("error in executing template %s", err)
		return err
	}
	log.Println("auth page gotten")
	return nil

}

func (app *Config) renderSign(w http.ResponseWriter, r string) error {
	render := "./cmd/web/templates/signup.page.go.html"
	t, err := template.ParseFiles(render)
	if err != nil {
		log.Println(err)
		return err
	}
	entry := SignUp{
		FirstName: "",
		LastName:  "",
		Email:     "",
		Password:  "",
	}

	err = t.ExecuteTemplate(w, "signup", entry)
	if err != nil {
		log.Println("error in executing template", err)
		return err
	}
	log.Println("SignUp page gotten")
	return nil

}
