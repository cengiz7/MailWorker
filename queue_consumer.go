package main

import (
	"github.com/streadway/amqp"
	"log"
)

func ConsumeFromQueue(){
	que_args["x-max-priority"] = priorityRange

	ch := OpenAmqpChannel()
	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		priorityQueueName, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		que_args,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	go StartConsume(msgs)
}


// using <-chan rather that chan amqp... is important
// |<- means| read only channel
func StartConsume(msgs <-chan amqp.Delivery){
	for d := range msgs {
		if current_driver == "sendgrid"{
			log.Printf(string(d.Body)+ "\n\n\n")
			d.Ack(false)
		}

	}
}