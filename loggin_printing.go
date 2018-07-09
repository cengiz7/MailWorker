package main

import (
	"log"
	"os"
)

// if you want  to append new logs to the existing log file, change the secon os.O_CREATE with os.RWONLY
func CreateLogFile(){
	f, err := os.OpenFile(logPath, os.O_RDWR | os.O_CREATE | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}


// if you don't want to log succesfull processes, then leave the sccsmsg field while calling the FailOnError func.
func LoggingChecking(err error, errmsg string,sccsmsg string) bool {
	if err != nil {
		log.Fatalf("%s: %s", errmsg, err)
		return true
	} else {
		if sccsmsg != ""{
			log.Printf("%s",sccsmsg)
		}
		return false
	}
}
