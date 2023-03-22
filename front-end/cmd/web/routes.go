package main

import (
	"net/http"
)

func (app *Config) routes() http.Handler {

	mux := http.NewServeMux() //an object to implement the router interface

	//gets the static files folder which contains the image
	fileServer := http.FileServer(http.Dir("./cmd/web/templates/static/")) //.gets to the root of the application
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")

		app.renderHome(w, "test.page.go.html")
	})

	return mux
}
