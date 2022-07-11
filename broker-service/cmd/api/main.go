package main

import (
	"fmt"
	"log"
	"net/http"
)

// broker service connects to the frontend,
//it processes each request sent by the frontend and sends a response back to it

const webport = "80"

type Config struct{}

func main() {
	app := Config{}
	log.Printf("starting broker service on port %s\n", webport)
	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}
	//start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
