package main

import (
	"net/http"
)

func (app *Config) routes() http.Handler {

	mux := http.NewServeMux() //an object to implement the router interface

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

		app.renderHome(w, "test.page.go.html")
	})

	return mux
}
