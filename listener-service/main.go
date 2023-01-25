package main

import (
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//connect to rabbit mq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMq")

	//start listening for messages

	//create a consumer to consume messages from the queque

	//watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	var counts int64

	//using backoff to attempt to connect a fixed number of times because
	//rabbbit mq is slow to connect
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Println("RabbitMq not yet ready.....")
			counts++
		} else {
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
