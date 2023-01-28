package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declares an exchange on the server.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name of the exchange
		"topic",      //type
		true,         //durable
		false,        //auto deleted
		false,        //not just used internally
		false,        //no wait
		nil,          //no specific arguments
	)
}

// declare a queue to hold messages and deliver to consumers
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //delete ehrn unused
		true,  //exclusive
		false, //no-wait
		nil,   //no specific arguments

	)
}
