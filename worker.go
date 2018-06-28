package main

import (
	"log"

	"github.com/streadway/amqp"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	args := make(amqp.Table)
	args["x-max-priority"] = int64(2)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"json_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")


	q, err := ch.QueueDeclare(
		"", // name
		true,   // durable
		false,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		args,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s",
		q.Name, "json_direct", "priority")
	err = ch.QueueBind(
		q.Name,        // queue name
		"priority",             // routing key
		"json_direct", // exchange
		false,
		args)
	failOnError(err, "Failed to bind a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		args,    // args
	)
	failOnError(err, "Failed to register a consumer")

	defer ch.Close()
	defer conn.Close()

	forever := make(chan bool)
	t := time.Duration(50)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(t * time.Millisecond)
			d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}