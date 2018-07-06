package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"flag"
	"github.com/streadway/amqp"
	"encoding/json"
	"strconv"
)

var (
	flagPort = flag.String("port", "9000", "Port to listen on")
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

// handles http post requests and send them directly to priority q
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		SendToPriorityQueue(&w,body)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ParseAndPublish(yy map[string]interface{}) error {
	for m := 0 ; m < len(yy["mails"].(I)); m++{
		body,err := json.Marshal(yy["mails"].(I)[m].(K)["mail"].(K))
		if err != nil {return err}
		p , err := strconv.Atoi(yy["mails"].(I)[m].(K)["priority"].(string))
		if err != nil {return err}
		// if priority value bigger than current priority range,
		// then we can say this is a timestamp
		if int64(p) > priorityRange {
			err = SaveToDb(p,body)
			fmt.Println("Veri tabanına kaydedildi")
		}else {
			err = PublishOverChannel(body,uint8(p))
		}
		if err != nil {return err}
	}
	return nil
}

func SendToPriorityQueue(w *http.ResponseWriter, body []byte){
	err := json.Unmarshal(body, &y)
	FailOnError(err,"Unmarshalling error before sending message to priority queue section.")
	if err == nil {
		err = ParseAndPublish(y)
	}
	RespondPostRequests(w,err)
}

func PublishOverChannel(body []byte,priority uint8) error {
	ch := OpenAmqpChannel()
	err := ch.Publish(
		"",     // exchange
		priorityQueueName, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            body ,
			DeliveryMode:    amqp.Persistent,
			Priority:        priority,
		})
	return err
}

func RespondPostRequests(w *http.ResponseWriter,err error){
	if err != nil {
		fmt.Fprint(*w, "Error : This message has not been published to queue.")
		log.Fatalf("Publishing error : %s", err)
	}else{
		fmt.Fprint(*w, "Message has been published to queue.")
	}
}


func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Custom-Header", "myvalue")
	w.Header().Set("Content-Type", "application/json")
	w.Write(nil)
}
func HttpWorker() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/", GetHandler)
	mux.HandleFunc("/post",PostHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}