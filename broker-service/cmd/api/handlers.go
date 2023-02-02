package main

import (
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"` //info needed to authenticate
	Log    LogPayload  `json:"log,omitempty"`  //info needed by to show a user is logged in
	//LogGRPc LogPayload  `json:"loggrpc,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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

// this handler handles all request to display in the front-end and called by the javascript
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	//taking action on the contents being received
	switch requestPayload.Action {
	//from authentication-service to broker
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	//the mail service to broker
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	//from the Logger	to broker
	case "log":
		app.LogItemViaRPC(w, requestPayload.Log)

	//from the Rabbit-Mq to broker
	case "queue":
		app.LogEventViaRabit(w, requestPayload.Log)

	default:
		app.ErrorJSON(w, errors.New("unknown Action"))
	}
}
