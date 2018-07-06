package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"
	"bytes"
)

var url = "http://localhost:9000/post"


func PrintResponseStatus(stat string,head http.Header,body string){
	fmt.Println("response Status:", stat)
	fmt.Println("response Headers:", head)
	//fmt.Println("response Body:", body)
}

func MakeRequest(body []byte){
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	respBody,_ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	PrintResponseStatus(resp.Status,resp.Header,string(respBody))
}

func SendAllMails(){
	fmt.Println("URL:>", url)
	mail, err := ioutil.ReadFile("mail_ornekleri.txt")
	if err != nil{
		log.Fatalf("Dosya acilamadi.")
		os.Exit(1)
	}
	MakeRequest(mail)
}