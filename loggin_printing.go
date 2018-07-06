package main

import (
	"log"
	"os"
)

func CreateLogFile(){
	f, err := os.OpenFile("logs.log", os.O_RDWR | os.O_CREATE | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
