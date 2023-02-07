package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"

	"log"
	"net/http"
)

func (app *Config) renderAuth(w http.Request, r *http.Request) (AuthPayload, error) {
	render := ".templates/auth.gohtml"
	t, err := template.New("api-html").ParseFiles(render)
	if err != nil {
		log.Println(err)
		return AuthPayload{}, err
	}
	var tab bytes.Buffer

	err = r.ParseForm()
	if err != nil {
		return AuthPayload{}, err
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	data := AuthPayload{
		Email:    email,
		Password: password,
	}
	err = t.ExecuteTemplate(&tab, "signup", data)
	if err != nil {
		return data, err
	}
	return data, nil

}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some json to send to the Authentication Microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the Authentication service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) //must be the same as the docker compose yml
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("Error in Calling Auth Service"))
		return
	}
	//create a variable We'll read response.body into
	var jsonFromService jsonResponse

	//decode the json from the Auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		log.Println(err)
		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	///writes the actual data from the service
	app.Writejson(w, http.StatusAccepted, payload)

}
