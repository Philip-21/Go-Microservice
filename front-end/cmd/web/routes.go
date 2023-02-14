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

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		app.renderAuth(w, "auth.page.go.html")
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		app.renderSign(w, "signup.page.go.html")

	})

	return mux
}
