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

// create a new consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		log.Println(err)
		return Consumer{}, err
	}
	return consumer, nil
}

// declare an amqp server
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

// communicates with the logger service when a listener send a reply
// to the broker
func LogEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	//an http client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()
	//if status is accepted
	if response.StatusCode != http.StatusAccepted {
		log.Println(err)
		return err
	}
	return nil
}

// handle what is sent
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log whatever we get
		//calls the url from logger service(WriteLog())
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
	//declare a queue to hold messages and deliver to consumers
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
	/*Consume immediately starts delivering queued messages.
	Begin receiving on the returned chan Delivery before any other operation on the Connection or Channel.
	Continues deliveries to the returned chan Delivery until Channel.Cancel, Connection.Close, Channel.Close, or an AMQP exception occurs.
	*/
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
			//the request fo
			go handlePayload(payload)
		}
	}()
	fmt.Printf("waiting or message [Exchange, Queue] [logs, topic, %s]\n", q.Name)
	<-forever
	return nil
}
