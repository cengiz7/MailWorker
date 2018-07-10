package main

import (
	"log"
	"os"
)

func CreateLogFile(){
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

// ..LOG LEVELS..
//----------------
// debug > level 0
// info  > level 1
// error > level 2
// fatal > level 3
//
func LoggingChecking(err error, errmsg string,sccsmsg string,level uint8) bool {
	if level >= logginLevel {
		if err != nil {
			log.Fatalf("%s: %s", errmsg, err)
			return true
		} else {
			if sccsmsg != "" && (logginLevel <= 1 ) {
				log.Printf("%s",sccsmsg)
			}
			return false
		}
	}else {
		if err !=nil {return true}
		return false
	}
}
