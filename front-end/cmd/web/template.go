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

	data := AuthPayload{
		Email:    "",
		Password: "",
	}
	render := ".templates/auth.gohtml"
	t, err := template.ParseFiles(render)
	if err != nil {
		log.Println(err)
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil

}
