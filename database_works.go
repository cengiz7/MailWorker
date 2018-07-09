package main

import (
	_ "github.com/mattn/go-sqlite3"
	"time"
	"fmt"
	"database/sql"
	"log"
	"os"
)

var database *sql.DB


func CreateDbTables(){
	db, err := sql.Open("sqlite3", dbPath)
	LoggingChecking(err,"While opening sqlite3 .db","Db created/opened.")
	database = db
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT)")
	LoggingChecking(err,"Wile creating table 'mails'","DB / mails table created.(if not exists)")
	statement.Exec()
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS troubled_mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT, error TEXT)")
	LoggingChecking(err,"Wile creating table 'troubled_mails'","DB / troubled_mails table created.(if not exists)")
	statement.Exec()
}

func SaveToDb(tmstmp int,body []byte) error {
	statement, err := database.Prepare("INSERT INTO mails (timestamps, body,tries) VALUES (?, ?,?)")
	statement.Exec(tmstmp, body,0)
	return err
}

func ConsumeFromDb(){
	var id uint64
	var timestamp int64
	var body []byte
	var ids []uint64
	for ;0<1;{
		ids = nil
		tmstmp := time.Now().Unix()
		rows, err := database.Query(fmt.Sprintf("SELECT id, timestamps , body FROM mails WHERE timestamps BETWEEN (%d) and (%d) or timestamps < (%d)",tmstmp-int64(dbQueryCheckRange),tmstmp,tmstmp))
		LoggingChecking(err,"db query error while consuming from db",fmt.Sprintf("Database check triggered at > %d",tmstmp))
		start := time.Now().Unix()
		for rows.Next(){
			rows.Scan(&id, &timestamp, &body)
			err = PublishOverChannel(body,3)
			LoggingChecking(err,"While publishing to queue at ConsumeFromDb",fmt.Sprintf("Message (No:%d) published to queue.",id))
			if err == nil {
				ids = append(ids,id)
			}
		}
		if ids != nil {
			err = DeleteFromDb(ids,"mails")
			LoggingChecking(err,"Deletefromdb returned error","")
		}
		finish := time.Now().Unix()   // if pusblishing took too much time, then abstract time from time.sleep duration
		time.Sleep(time.Duration(dbQueryCheckRange - diff(start,finish) ) * time.Second)
	}
}

func DeleteFromDb(ids []uint64,where string) error {
	for i := range ids{
		stmt, err := database.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = %d",where,ids[i]))
		LoggingChecking(err,"While deleting from mails",fmt.Sprintf("Message (No:%d) deleted from DB/%s.",ids[i],where))
		stmt.Exec()
		if err != nil{return err}
	}
	return nil
}

func diff(strt,fnsh int64) int{
	if dbQueryCheckRange < dbCheckPeriot {
		log.Fatalf("dbQueryCheckRange can't be Equal or Lower than query check range!\n Fix it under config.txt file. ")
		os.Exit(1)
	}
	if fnsh-strt > int64(dbQueryCheckRange- dbCheckPeriot) {
		additionTime := fnsh-strt-int64(dbCheckPeriot-dbQueryCheckRange)
		return int(additionTime)
	} else {
		return 0
	}
}