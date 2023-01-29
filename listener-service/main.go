package main

import (
	"listener-service/event"
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

	//start listening for messages
	log.Println("listenin for and  consuming RabbitMQ messages...")

	//create a consumer to consume messages from the queque
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		panic(err)
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
