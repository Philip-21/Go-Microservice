package main

import (
	"log"
	"net/http"
	"net/rpc"
)

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
