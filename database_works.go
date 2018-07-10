package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"os"
	"time"
	"fmt"
	"log"
)

var database *sql.DB


func CreateDbTables(){
	db, err := sql.Open("sqlite3", dbPath)
	if LoggingChecking(err,"While opening sqlite3 .db","Db created/opened.",3){
		os.Exit(7)
	}
	database = db
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT)")
	if LoggingChecking(err,"Wile creating table 'mails'","DB / mails table created.(if not exists)",2) {
		os.Exit(8)
	}
	statement.Exec()
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS troubled_mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT, error TEXT)")
	if LoggingChecking(err,"Wile creating table 'troubled_mails'","DB / troubled_mails table created.(if not exists)",2){
		os.Exit(9)
	}
	statement.Exec()
}

func SaveToDb(tmstmp int,body []byte) error {
	statement, err := database.Prepare("INSERT INTO mails (timestamps, body,tries) VALUES (?, ?,?)")
	LoggingChecking(err,"While inserting data into 'mails' table","Mail data successfully inserted into db.",2)
	statement.Exec(tmstmp, body,0)
	return err   // it reteturns error through http request responding function
}

/* checking pending mails from mails table at each given dbCheckPeriot seconds time.
 if publishing and deleting takes too much, then it calculates delay time and eliminate
 delay timing. So, no messages from table will be missed.
*/
func ConsumeFromDb(){
	var id uint64
	var timestamp int64
	var body []byte
	var ids []uint64
	for ;0<1;{
		ids = nil
		tmstmp := time.Now().Unix()
		rows, err := database.Query(fmt.Sprintf("SELECT id, timestamps , body FROM mails WHERE timestamps BETWEEN (%d) and (%d) or timestamps < (%d)",tmstmp-int64(dbQueryCheckRange),tmstmp,tmstmp))
		LoggingChecking(err,"db query error while consuming from db",fmt.Sprintf("Database check triggered at > %d",tmstmp),2)
		start := time.Now().Unix()
		for rows.Next(){
			rows.Scan(&id, &timestamp, &body)
			err = PublishOverChannel(body,3)
			LoggingChecking(err,"While publishing to queue at ConsumeFromDb",fmt.Sprintf("Message (No:%d) published to queue.",id),2)
			if err == nil {
				ids = append(ids,id)
			}
		}
		if ids != nil {
			err = DeleteFromDb(ids,"mails")
			LoggingChecking(err,"Deletefromdb returned error","Messages successfully deleted from 'mails' table",2)
		}
		finish := time.Now().Unix()   // if pusblishing took too much time, then abstract time from time.sleep duration
		time.Sleep(time.Duration(dbQueryCheckRange - diff(start,finish) ) * time.Second)
	}
}

// deleting messages after consuming.
func DeleteFromDb(ids []uint64,where string) error {
	for i := range ids{
		stmt, err := database.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = %d",where,ids[i]))
		LoggingChecking(err,"While deleting message from "+where+ " table.",fmt.Sprintf("Message (No:%d) deleted from DB/%s.",ids[i],where),2)
		stmt.Exec()
		if err != nil{return err}
	}
	return nil
}

// calculating delay time and abstract it from dbCheckPeriot time
func diff(strt,fnsh int64) int{
	if dbQueryCheckRange < dbCheckPeriot {
		log.Fatalf("dbQueryCheckRange can't be Equal or Lower than query check range!\n Fix it under config.txt file. ")
		os.Exit(10)
	}
	if fnsh-strt > int64(dbQueryCheckRange- dbCheckPeriot) {
		additionTime := fnsh-strt-int64(dbCheckPeriot-dbQueryCheckRange)
		return int(additionTime)
	} else {
		return 0
	}
}