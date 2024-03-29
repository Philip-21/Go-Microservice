package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// communicates with the logger service when a User
// sends a request and saves the request in the db
func LogEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {

		return err
	}
	request.Header.Set("Content-Type", "application/json")
	//an http client
	client := &http.Client{}
	//Do sends an HTTP request and returns an HTTP response,
	// following policy (such as redirects, cookies, auth)
	//as configured on the client.
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()
	//if status is accepted
	if response.StatusCode != http.StatusAccepted {

		return err
	}
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

// listens to the queue
func (consumer *Consumer) Listen(topics []string) error {
	//Channel opens a unique concurrent server channel,
	// to process the bulk of AMQP messages.
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
