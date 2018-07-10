package main

import (
	"github.com/streadway/amqp"
)

func ConsumeFromQueue(){
	que_args["x-max-priority"] = priorityRange

	ch := OpenAmqpChannel()
	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	LoggingChecking(err, "Failed to set QoS","QoS successfully set for "+priorityQueueName,2)

	msgs, err := ch.Consume(
		priorityQueueName, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		que_args,    // args
	)
	LoggingChecking(err, "Failed to register a consumer","Consumer successfully set for "+priorityQueueName,2)

	go StartConsume(msgs)
}


// using <-chan rather that chan amqp... is important
// |<- means| read only channel
func StartConsume(msgs <-chan amqp.Delivery){
	LoggingChecking(nil,"","Queue consuming started..",1)
	for d := range msgs {
		if currentDriver == "sendgrid"{
			LoggingChecking(nil,"",string(d.Body[:55])+ "\n\n\n",1)
			d.Ack(false)
		}

	}
}