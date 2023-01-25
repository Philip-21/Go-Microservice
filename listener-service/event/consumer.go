package event

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// receiving events from the queue
type Consumer struct {
	conn      *amqp.Connection
	queuename string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

func (consumer *Consumer) setup() error {
	//Channel opens a unique, concurrent server channel to
	// process the bulk of AMQP messages
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// listens to the queue
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	//Close initiate a clean channel closure by sending a close
	// message with the error code set to '200'
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	//range through topics
	for _, s := range topics {
		//bind s to a queue
		ch.QueueBind(
			q.Name,
			s,
			"log_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil
	}

	//consume everything that comes from rabbitmq until application exits
	//declare a channel
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Printf("waiting or message [Exchange, Queue] [logs, topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log whatever we get
		err := LogEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//authenticate

		//a default case that logs things with information thathas some value
	default:
		err := LogEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func LogEvent(payload Payload) error {
	return nil
}
