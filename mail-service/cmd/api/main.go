package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service n webport %s", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
