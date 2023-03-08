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
	//login
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	//the mail service to broker
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	/*
		from the Logger to the broker based on Api Request
		either grpc or Rest
		--------------For grpc----------------,
		the grpc handles the same action as REST Based on  Protocol buffers
		for transmitting requests
		for Grpc , the Reqeust is Called based on the button clicked ,
		based on the requestPayload ,the RestApi handles the reqeust  from the client then
		calls the grpc server to execute the request and send a response back to the to the user.
		it calls the LogGRpc function which does the grpc operations to handle the request and submit a response
		which is sent back in json format to the user
	*/
	case "log":
		app.logItem(w, requestPayload.Log)
	//handles request from the client using RestApi , calls the
	//the Rpc to handle reuest and returns a response back in json
	case "rpc":
		app.LogItemViaRPC(w, requestPayload.LogRPc)
	case "queue":
		app.LogEventViaRabit(w, requestPayload.Log)

	default:
		app.ErrorJSON(w, errors.New("unknown Action"))
	}
}
