package main

import (
	"time"
)

// globals
var current_driver = "sendgrid"   // required while consuming messages from queue
var priorityQueueName = "priority_queue"
var priorityRange = int64(3)  // 0 lowest 3 is highest pri.
// dbcheckperiot should not be higher than checkrange
var dbCheckPeriot = 30   // seconds
var dbQueryCheckRange = 35 // seconds
var y map[string]interface{}
type I = []interface {}
type K = map[string]interface {}


func main(){
	CreateLogFile()
	CreateDbTables()
	DeclarePriorityQueue()
	go HttpWorker()
	time.Sleep(time.Duration(1)*time.Second)
	go ConsumeFromDb()
	ConsumeFromQueue()
	for i := 0 ; i< 10; i++{
		SendAllMails()
	}
	ch := make(chan bool)
	<-ch
}