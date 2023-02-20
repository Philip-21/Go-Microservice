package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// broker service connects to the frontend,
//it processes each request sent by the frontend and sends a response back to it

const webport = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("starting broker service on port %s\n", webport)
	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}
	//start server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64

	//using backoff to attempt to connect a fixed number of times because
	//rabbbit mq is slow to connect
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq") //speaks to the rabbit-mq docker compose service
		if err != nil {
			log.Println("RabbitMq not yet ready.....")
			counts++
		} else {
			log.Println("Connected to RabbitMq")
			connection = c
			break
		}
		if counts > 5 {
			log.Println(err)
			return nil, err
		}
		//increase the delay
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second //raising it to the power of 2
		log.Println("Backing off....")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
