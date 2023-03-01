package main

import (
	"authentication/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models database.Models
}

func main() {

	log.Println("Starting authentication service")
	// connect to db
	conn := database.ConnectToDB()
	if conn == nil {
		log.Panic("Cant connect to Postgres!")
	}
	log.Println("connected to db ")

	// database.SeedDB(conn)
	// log.Println("seeding db completed")

	// set up config
	app := Config{
		DB:     conn,
		Models: database.New(conn),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
