package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// communictes with the logger service when someone succesfully authenticates
// and sends the auth details  into mongoDb
func (app *Config) logAuth(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil

}

func (app *Config) LogSignUp(name string, firstname string, lastname string, email string) error {
	var dbentry struct {
		Email     string `json:"email"`
		EntryName string `json:"entryname"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}
	dbentry.EntryName = name
	dbentry.Email = email
	dbentry.FirstName = firstname
	dbentry.LastName = lastname

	Json, _ := json.MarshalIndent(dbentry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(Json))
	if err != nil {
		log.Println(err)
		return err
	}
	//an http client to transport the request to the logger service
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (app *Config) SignUp(w http.ResponseWriter, r *http.Request) {
	var SignUp struct {
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Password  string `json:"password"`
	}

	err := app.ReadJSON(w, r, &SignUp)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(SignUp.Password), 8)
	//svae into postres database
	user, err := app.Models.User.CreateUser(SignUp.FirstName, SignUp.LastName, SignUp.Email, string(hashPassword))
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}
	//call the logger service and save to the mongo database
	err = app.LogSignUp("signup", user.Email, user.FirstName, user.LastName)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("SignUp user %s", user.FirstName),
		Data:    user,
	}
	app.Writejson(w, http.StatusAccepted, payload)

}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	//declare a variable that haa the same tags with the json
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	//validate the user from the db
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
	}
	// log authentication and send the logged details to logger service
	err = app.logAuth("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in User %s", user.Email),
		Data:    user,
	}
	app.Writejson(w, http.StatusAccepted, payload)
}
