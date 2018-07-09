package main

import (
	"log"
	"bufio"
	"strings"
	"regexp"
	"os"
	"strconv"
	"fmt"
)

// Globals
// for type shortening > use (I) rather than []interface {}
type I = []interface {}
type K = map[string]interface {}

var (
	configFile = "includes/config.txt"
	currentDriver string   // required while consuming messages from queue
	priorityQueueName string
	priorityRange int64  // 0 lowest 3 is highest pri.
	// dbcheckperiot should not be higher than checkrange
	dbCheckPeriot int   // seconds
	dbQueryCheckRange int // seconds
	dbPath string
	logPath string
	y map[string]interface{}
	config map[string]string
)

func main(){
	DeclarePriorityQueue()
	CreateDbTables()
	go HttpWorker()
	go ConsumeFromDb()
	ConsumeFromQueue()

	for i := 0 ; i< 10; i++{
		SendAllMails()
	}
	ch := make(chan bool)
	<-ch
}

func init(){
	// parse includes/config.txt into a map
	//
	// format of includes/config.txt:
	//   key=value
	//
	// access by:
	//  value := config["key"]
	//
	config = ReadConfig(configFile)
	dbPath = config["dbPath"]
	currentDriver = config["currentDriver"]
	priorityQueueName = config["priorityQueueName"]
	logPath = config["logPath"]
	var err error
	priorityRange, err = strconv.ParseInt(config["priorityRange"], 10, 64)
	LoggingChecking(err,"Error while parsing priorityRange from config.txt","")
	dbCheckPeriot,err = strconv.Atoi(config["dbCheckPeriot"])
	LoggingChecking(err,"Error while parsing dbCheckPeriot from config.txt","")
	dbQueryCheckRange, err = strconv.Atoi(config["dbQueryCheckRange"])
	LoggingChecking(err,"Error while parsing dbQueryCheckRange from config.txt","")

	CreateLogFile()
}

func ReadConfig(filename_fullpath string) map[string]string {
	prg := "ReadConfig()"
	file, err := os.Open(filename_fullpath)
	if LoggingChecking(err,fmt.Sprintf("%s: os.Open(): %s\n", prg, err)," Config file opened."){
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	options := ScanFile(scanner)
	return options
}

func ScanFile(scanner *bufio.Scanner) (map[string]string) {
	prg := "ReadConfig()"
	var options map[string]string
	options = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") == true {
			re, err := regexp.Compile(`([^=]+)=(.*)`)
			if LoggingChecking(err,fmt.Sprintf("%s: regexp.Compile(): error=%s",prg, err),"") {
				os.Exit(1)
			} else {
				confOption := re.FindStringSubmatch(line)[1]
				confValue := re.FindStringSubmatch(line)[2]
				options[confOption] = confValue
				log.Printf("%s: out[]: %s ... confOption=%s, confValue=%s\n",prg,line,confOption,confValue)
			}
		}
	}
	log.Printf("%s: options[]: %+v\n", prg, options)
	if LoggingChecking(scanner.Err(),fmt.Sprintf("%s: scanner.Err(): %s\n", prg, scanner.Err()),""){
		os.Exit(1)
	}
	return options
}
