package main

import (
	"broker/logs"
	"context"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RPCPayload struct {
	Name string
	Data string
}

// Send requests using RPC and save into logger-Service database then displays result in the frontend
func (app *Config) LogItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001") //calls from the docker compose file
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

// sends request using GRPC connects through this broker-sevice
// then save's to the logger-service database and displays result on the front-end
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	//getting the request from the user
	var requestPayload RequestPayload

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	///connect to GRPC  server
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	defer conn.Close()
	//create a client
	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//Call the writeLog()
	//write to the log using GRPC
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	//send a response back
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"
	app.Writejson(w, http.StatusAccepted, payload)
}
