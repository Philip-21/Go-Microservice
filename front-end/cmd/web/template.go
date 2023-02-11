package main

import (
	"html/template"
	"log"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Config struct{}

func renderAuth(w http.ResponseWriter, v string) error {

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
		return err
	}
	return nil

}
