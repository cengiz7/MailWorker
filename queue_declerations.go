package main

import "github.com/streadway/amqp"

var que_args = make(amqp.Table)


func ConnectAmqp()(*amqp.Connection){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func OpenAmqpChannel()(*amqp.Channel){
	conn := ConnectAmqp()
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	return ch

}

func DeclarePriorityQueue(){
	que_args["x-max-priority"] = priorityRange
	ch := OpenAmqpChannel()
	_, err := ch.QueueDeclare(
		priorityQueueName, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		que_args,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")
	ch.Close()
}