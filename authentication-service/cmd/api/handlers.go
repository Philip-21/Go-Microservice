package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	log.Println("Saved  details into Mongo service ")
	return nil
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
		log.Println("Error in reading JSON")
		return
	}
	//validate the user from the db
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		log.Println("Invalid Email")
		return
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		log.Println("Invalid password")
	}
	// log authentication and send the logged details to logger service
	err = app.logAuth("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.ErrorJSON(w, err)
		log.Println(err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in User %s", user.Email),
		Data:    user,
	}
	app.Writejson(w, http.StatusAccepted, payload)
	log.Println("Authenticated")
}
