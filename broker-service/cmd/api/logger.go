package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/rpc"
)

// The broker can send  requests using Api's and save it on the logger
// servcie database then displays it in the frontend
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

// to handle logging and pushing events to queue
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

type RPCPayload struct {
	Name string
	Data string
}

// Send requets using RPC and save into logger-Service database then displays result in the frontend
func (app *Config) LogItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	rpcpayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	//get a result call
	var result string
	//call the rpc server
	err = client.Call("RPCServer.LogInfo", rpcpayload, &result)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.Writejson(w, http.StatusAccepted, payload)
}
