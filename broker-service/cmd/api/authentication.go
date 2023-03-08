package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some json to send to the Authentication Microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the Authentication service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData)) //must be the same as the docker compose yml
	if err != nil {
		app.ErrorJSON(w, err)
		log.Println("Error in writing json")
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		log.Println("Error in creating client")
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid credentials"))
		log.Println("Invalid Credentials")
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("Error in Calling Auth Service"))
		log.Println("Error in Calling Auth Service")
		return
	}

	//create a variable We'll read response.body into
	var jsonFromService jsonResponse

	//decode the json from the Auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.ErrorJSON(w, err)
		log.Println("error in decoding")
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
	log.Println("Authenticated by Broker Service")

	///writes the actual data from the service
	app.Writejson(w, http.StatusAccepted, payload)

}
