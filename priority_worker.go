package main

import (
	"log"

	"github.com/streadway/amqp"
	"time"
	"os"
)



func priorityWorker() {
	f, err := os.OpenFile("Queue_logs.txt", os.O_RDWR | os.O_CREATE | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	args := make(amqp.Table)
	args["x-max-priority"] = int64(2)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")


	q, err := ch.QueueDeclare(
		"priority_queue", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		args,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

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
			log.Printf("\n\nPriority queue >> Received a message: %s ...\n", d.Body[0:50])
			time.Sleep(t * time.Millisecond)
			log.Printf("Done\n")
			d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}