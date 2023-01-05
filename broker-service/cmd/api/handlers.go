package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"` //info needed to authenticate
	Log    LogPayload  `json:"log,omitempty"`  //info needed by to show a user is logged in
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"email"`
	Data string `json:"data"`
}

// the output diplayed in the frontend
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` //used any data type  instead of an interface , cause im parsin  insmall values
}

// the homepage handler for the frontend
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.Writejson(w, http.StatusOK, payload)
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some json to send to the Authentication Microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the Authentication service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) //must be the same as the docker compose yml
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
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
		app.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
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
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	//an http client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()
	//if status is accepted
	if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, err)
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	app.Writejson(w, http.StatusAccepted, payload)
}

// this handler handles all request
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	//taking action on the contents being received
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		app.ErrorJSON(w, errors.New("unknown Action"))
	}
}
