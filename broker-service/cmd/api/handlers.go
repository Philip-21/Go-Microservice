package main

import (
	"errors"
	"net/http"
)

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
		app.authenticate(w, r)
	//the mail service to broker
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	//from the Logger	to broker
	case "log":
		app.logItem(w, requestPayload.Log)

	//from the Rabbit-Mq to broker
	case "queue":
		app.LogEventViaRabit(w, requestPayload.Log)

	default:
		app.ErrorJSON(w, errors.New("unknown Action"))
	}
}
