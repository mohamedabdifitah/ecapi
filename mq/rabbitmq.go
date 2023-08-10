package mq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Ch *amqp.Channel
)

func ConnectMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("rabbitMQ connection established")
	Ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer Ch.Close()
	_, err = Ch.QueueDeclare(
		"new order", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

}
func PublishOrder[T interface{}](order T, queue string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("rabbitMQ connection established")
	Ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer Ch.Close()
	_, err = Ch.QueueDeclare(
		"new order", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	encoder := new(bytes.Buffer)
	json.NewEncoder(encoder).Encode(order)
	err = Ch.PublishWithContext(ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encoder.Bytes(),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent order ")
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
