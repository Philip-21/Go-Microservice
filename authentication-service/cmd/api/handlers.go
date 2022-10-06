package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authentication(w http.ResponseWriter, r *http.Request) {

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

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in User %s", user.Email),
		Data:    user,
	}
	app.Writejson(w, http.StatusAccepted, payload)
}
