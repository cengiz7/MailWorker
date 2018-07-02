package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	url := "http://localhost:9000/post"
	fmt.Println("URL:>", url)
	mail, err := ioutil.ReadFile("mail_ornekleri.txt")
	if err != nil{
		log.Fatalf("Dosya acilamadi.")
		os.Exit(1)
	}
	mailstr := string(mail)
	mailler := strings.Split(mailstr,"sep-from-here")
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

}