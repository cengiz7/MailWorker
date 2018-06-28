package main

import (
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
if err != nil {
log.Fatalf("%s: %s", msg, err)
}
}

var deneme = []string{
	"low","low","high","low","medium","low","low","high","high","medium",
	"low"}
var deneme1 = []uint8 {
	0,0,2,0,1,0,0,2,2,1,0}

func main() {
	args := make(amqp.Table)
	args["x-max-priority"] = int64(2)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"json_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		args,           // arguments
	)

	failOnError(err, "Failed to declare a queue")

	// for testing !!!!!
	body :=""
	for k := 0; k<50 ;k++{
		for i:=0;i<11;i++{
			body = deneme[i]
			err = ch.Publish(
				"json_direct",     // exchange
				"priority", // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "text/plain",
					ContentEncoding: "",
					Body:            []byte(body),
					DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
					Priority:        deneme1[i],
				})
			//log.Printf(" [x] Sent %s", body)
			failOnError(err, "Failed to publish a message")
		}

	}


}
