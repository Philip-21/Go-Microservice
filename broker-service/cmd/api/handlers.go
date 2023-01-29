package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"` //info needed to authenticate
	Log    LogPayload  `json:"log,omitempty"`  //info needed by to show a user is logged in
	Mail   MailPayload `json:"mail,omitempty"`
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

/*func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
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
}*/

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

// to handle logging and emmit an event to rabbit mq
func (app *Config) LogEventViaRabit(w http.ResponseWriter, l LogPayload) {
	err := app.PushToQueue(l.Name, l.Data)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	//send back a jsn response
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged Via RabbitMQ"

	app.Writejson(w, http.StatusAccepted, payload)
}
func (app *Config) PushToQueue(name, msg string) error {
	emmiter, err := event.NewEventEmmitter(app.Rabbit)
	if err != nil {
		log.Println(err)
		return err
	}
	//create a payload to store in queue'
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	j, _ := json.MarshalIndent(&payload, "", "\t")
	//push the payload to queue
	err = emmiter.Push(string(j), "Log.INFO")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
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
		app.LogEventViaRabit(w, requestPayload.Log)
	default:
		app.ErrorJSON(w, errors.New("unknown Action"))
	}
}
