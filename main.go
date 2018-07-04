package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)
var mailchan = make(chan uint)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	f, err := os.OpenFile("Queue_logs.txt", os.O_RDWR | os.O_CREATE | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	go func() {
		httpWorker()
	}()
	go func() {
		timeWorker()
	}()
	go func() {
		priorityWorker()
	}()

	time.Sleep(time.Duration(2) * time.Second)
	url := "http://localhost:9000/post"
	fmt.Println("URL:>", url)
	mail, err := ioutil.ReadFile("mail_ornekleri.txt")
	if err != nil{
		log.Fatalf("Dosya acilamadi.")
		os.Exit(1)
	}
	mailstr := string(mail)
	mailler := strings.Split(mailstr,"sep-from-here")
	forever := make(chan bool)
	for _,bb := range mailler{
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bb)))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
	<-forever
}