package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/rpc"
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

// Send requests using RPC and save into logger-Service database then displays result in the frontend
func (app *Config) LogItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001") //calls from the docker compose file
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}
	var r *http.Request
	err = r.ParseForm()
	if err != nil {
		log.Println("error in getting form")
		return
	}
	// name := r.Form.Get("name")
	// Data := r.Form.Get("data")
	rpcpayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}
	// rpc, err := app.renderRpc(w, r)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
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
