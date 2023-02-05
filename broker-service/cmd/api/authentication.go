package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
)

func (app *Config) renderRpc(w http.ResponseWriter, r *http.Request) (RPCPayload, error) {
	render := ".templates/rpc.gohtml"
	t, err := template.New("rpc-html").ParseFiles(render)
	if err != nil {
		log.Println(err)
		return RPCPayload{}, err
	}
	var tbf bytes.Buffer

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		log.Println("error in getting form")
		return RPCPayload{}, err
	}
	data := RPCPayload{
		Name: r.Form.Get("name"),
		Data: r.Form.Get("data"),
	}
	err = t.ExecuteTemplate(&tbf, "rpc", data)
	if err != nil {
		return data, err
	}
	return data, nil
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
