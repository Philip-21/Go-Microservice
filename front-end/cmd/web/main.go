package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":80"

type Config struct{}

func main() {

	fmt.Println("Starting front end service on port ", portNumber)
	log.Println("frontend started")
	//define http server
	app := Config{}
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s", portNumber),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
