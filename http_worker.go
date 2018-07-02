package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"github.com/streadway/amqp"
	"strconv"
)

var (
	flagPort = flag.String("port", "9000", "Port to listen on")
)

var conn, err1 = amqp.Dial("amqp://guest:guest@localhost:5672/")
var ch, err2 = conn.Channel()
var y map[string]interface{}
type I = []interface {}
type K = map[string]interface {}
var ornek []byte


func GetHandler(w http.ResponseWriter, r *http.Request) {
	/*jsonBody, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
	}*/
	w.Header().Set("X-Custom-Header", "myvalue")
	w.Header().Set("Content-Type", "application/json")

	w.Write(ornek)
}



func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		if err := json.Unmarshal(body, &y); err != nil {
			fmt.Println(err)
			return
		}

		for m :=0 ; m< len(y["mails"].(I)) ; m++{
			if len(y["mails"].(I)[m].(K)["priority"].(string)) > 1 {
				println("zaman öncelikli > ",y["mails"].(I)[m].(K)["priority"].(string))
				err = ch.Publish(
					"",           // exchange
					"time_queue",       // routing key
					false,        // mandatory
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "application/json",
						Body:         body,
					})
			}else {
				i ,_:= strconv.Atoi(y["mails"].(I)[m].(K)["priority"].(string))
				p := uint8(i)
				println(p)
				err = ch.Publish(
					"",     // exchange
					"priority_queue", // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						Headers:         amqp.Table{},
						ContentType:     "application/json",
						ContentEncoding: "",
						Body:            body ,
						DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
						Priority:        p,
					})
				if err != nil {
					log.Fatalf("Publishing error : %s", err)
					break
				}
			}
		}
		fmt.Fprint(w, "POST done\n")

		/*
		the_list := y["mails"].(I)["priority"]
		for n, v := range the_list {
			fmt.Printf("index:%d  \nvalue:%v  \nkind:%s  \ntype:%s \n", n, v, reflect.TypeOf(v).Kind(), reflect.TypeOf(v))
		}*/

		//fmt.Println(y["mails"].(I)[0].(K)["content"].(K)["attachments"].(I)[0])
		//fmt.Println(len(y["mails"].(I)[0].(K)["content"].(K)["attachments"].(I)))



		/*
		for m := 0 ; m < len(y["mails"].(I)); m++{
			
			contentHeaders := y["mails"].(I)[m].(K)["content"].(K)["headers"].(K)
			contentBody := y["mails"].(I)[m].(K)["content"].(K)["body"].(K)
			mailAdress := y["mails"].(I)[m].(K)["content"].(K)["headers"].(K)["To"].(I)



			for _,val := range mailAdress {
				fmt.Println("Priority :",y["mails"].(I)[m].(K)["priority"].(string))
				fmt.Println("From :",contentHeaders["From"].(string))
				fmt.Println("To :",val)
				fmt.Println("Subject :",contentHeaders["Subject"].(string))
				for bodyType,bodyValue := range contentBody{
					fmt.Println("content-type :",bodyType, " value :",bodyValue)
				}
				fmt.Println("\n\n")
		}}*/

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func httpWorker() {
	//results = append(results, time.Now().Format(time.RFC3339))

	mux := http.NewServeMux()
	//mux.HandleFunc("/", GetHandler)
	mux.HandleFunc("/post", PostHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}