package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueMessage struct {
	mime     string
	body     []byte
	exchange string
	key      string
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		return
	}
}

// var channel = make(chan *amqp.Channel)
var channel = make(chan QueueMessage)

func InitAmqp() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("connect to RabbitMQ")
	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	for message := range channel {
		err := PublishMessage(ch, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func ProduceMessage(exchange, key, mime string, body interface{}) {
	bytedata, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	if mime == "" {
		mime = "application/json"
	}
	message := QueueMessage{
		mime:     mime,
		body:     bytedata,
		exchange: exchange,
		key:      key,
	}
	channel <- message
}
func PublishMessage(ch *amqp.Channel, message QueueMessage) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := ch.PublishWithContext(
		ctx,
		message.exchange, // exchange
		message.key,      // routing key (queue name)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: message.mime,
			Body:        message.body,
		})
	if err != nil {
		return err
	}
	return nil
}
