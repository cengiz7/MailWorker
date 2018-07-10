package main

import "github.com/streadway/amqp"

var que_args = make(amqp.Table)


func ConnectAmqp()(*amqp.Connection){
	conn, err := amqp.Dial(rabbitMqConnection)
	LoggingChecking(err, "Failed to connect to RabbitMQ","Succesfully connected to RabbitMQ",3)
	return conn
}

func OpenAmqpChannel()(*amqp.Channel){
	conn := ConnectAmqp()
	ch, err := conn.Channel()
	LoggingChecking(err, "Failed to open a channel","Successfully opened a amqp channel.",3)
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
	LoggingChecking(err, "Failed to declare a queue","Successfully declare "+priorityQueueName+ "on RabbitMQ.",3)
	ch.Close()
}