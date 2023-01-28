//pushes events on the queue

package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emmitter struct {
	connection *amqp.Connection
}

func (e *Emmitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emmitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	defer channel.Close()
	log.Println("Pushing to channel")

	/*Publish sends a Publishing from the client to an exchange on the server.
	When you want a single message to be delivered to a single queue,
	you can publish to the default exchange with the
	routingKey of the queue name. This is because every declared queue gets an implicit route to the default exchange*/
	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// /creating a push event emmiter
func NewEventEmmitter(conn *amqp.Connection) (Emmitter, error) {
	emmiter := Emmitter{
		connection: conn,
	}
	err := emmiter.setup() //using receivers as inheritance here
	if err != nil {
		return Emmitter{}, err
	}
	return emmiter, nil
}
