package main

import (
	"github.com/streadway/amqp"
	"log"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"time"
)
var z map[string]interface{}
var h map[string]interface{}
var t = time.Duration(50)

func checkErr(err error,s string) {
	if err != nil {
		panic(err)
		log.Fatalf(s)
	}
}

func timeWorker() {
	database, _ := sql.Open("sqlite3", "./pending_mails.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS troubled_mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT, error TEXT)")
	statement.Exec()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"time_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var id uint64
	var timestamp int64
	var body []byte
	var jobs bool
	go func(){
		for ;0<1;{
			tmstmp := time.Now().Unix()
			fmt.Printf("\nSELECT id, timestamps , body FROM mails WHERE timestamps BETWEEN (%d) and (%d) or timestamps < (%d)",tmstmp-int64(35),tmstmp,tmstmp)
			rows, err := database.Query(fmt.Sprintf("SELECT id, timestamps , body FROM mails WHERE timestamps BETWEEN (%d) and (%d) or timestamps < (%d)",tmstmp-int64(35),tmstmp,tmstmp))
			checkErr(err,"select query error")

			for rows.Next() {
				jobs = true
				rows.Scan(&id, &timestamp, &body)
				fmt.Println("\n"+strconv.FormatUint(id,10) + ": " + string(timestamp) + " > " + string(body[:60]) + " ....\n")
			}
			if jobs {
				fmt.Println("\n >> Jobs finished.. Waiting for new jobs\n")
				jobs = false
			}else {
				fmt.Println("\n >> Nothind done at this cycle.\n")
			}

			stmt, err := database.Prepare("DELETE FROM mails WHERE timestamps BETWEEN (?) and (?) or timestamps < (?)")
			checkErr(err,"delete row error")
			stmt.Exec(tmstmp-int64(35),tmstmp,tmstmp)

			time.Sleep(30 * time.Second)
		}
	}()


	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := json.Unmarshal(d.Body, &z)
			checkErr(err,"Unmarshalling error")
			i ,_:= strconv.Atoi(z["mails"].(I)[0].(K)["priority"].(string))
			time.Sleep(t * time.Millisecond)
			statement, _ = database.Prepare("INSERT INTO mails (timestamps, body) VALUES (?, ?)")
			statement.Exec(i, d.Body)
			log.Printf("Done\n")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages from 'time_queue'. To exit press CTRL+C")
	<-forever
}