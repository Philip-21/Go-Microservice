package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// calls the mail microservice and sends the result  back to the user
// the mail service doesn't send details to the logger-service
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mail service
	mailServiceURL := "http://mail-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		log.Println(err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}
	// send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.Writejson(w, http.StatusAccepted, payload)

}
