package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":80"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderHome(w, "test.page.go.html")
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		renderAuth(w, "auth.page.go.html")
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		renderSign(w, "signup.page.go.html")
	})

	fmt.Println("Starting front end service on port ", portNumber)
	log.Println("frontend started")
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Panic(err)
	}
}
