package main

import "database/sql"
import (
	_ "github.com/mattn/go-sqlite3"
	"time"
	"fmt"
)

var database *sql.DB


func CreateDbTables(){
	db, err := sql.Open("sqlite3", "./pending_mails.db")
	FailOnError(err,"While opening sqlite3 .db")
	database = db
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT)")
	FailOnError(err,"Wile creating table 'mails'")
	statement.Exec()
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS troubled_mails (id INTEGER PRIMARY KEY , timestamps INTEGER , body BLOB, tries TINYINT, error TEXT)")
	FailOnError(err,"Wile creating table 'troubled_mails'")
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
		fmt.Println("Database Check start at > ",tmstmp)
		FailOnError(err,"db query error while consuming from db")
		start := time.Now().Unix()
		for rows.Next(){
			rows.Scan(&id, &timestamp, &body)
			err = PublishOverChannel(body,3)
			FailOnError(err,"While publishing to queue at ConsumeFromDb")
			if err == nil {
				ids = append(ids,id)
			}
		}
		if ids != nil {
			err = DeleteFromDb(ids,"mails")
			FailOnError(err,"Deletefromdb returned error")
		}
		finish := time.Now().Unix()   // if pusblishing took too much time, then abstract time from time.sleep duration
		time.Sleep(time.Duration(dbQueryCheckRange - diff(start,finish) ) * time.Second)
	}
}

func DeleteFromDb(ids []uint64,where string) error {
	for i := range ids{
		stmt, err := database.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = %d",where,ids[i]))
		FailOnError(err,"While deleting from mails")
		stmt.Exec()
		if err != nil{return err}
	}
	return nil
}

func diff(strt,fnsh int64) int{
	if fnsh-strt > int64(dbQueryCheckRange- dbCheckPeriot) {
		additionTime := fnsh-strt-int64(dbCheckPeriot-dbQueryCheckRange)
		return int(additionTime)
	} else {
		return 0
	}
}