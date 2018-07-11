package main

import (
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
	rabbitMqConnection string
	priorityQueueName string
	priorityRange int64  // 0 lowest 3 is highest pri.
	// dbcheckperiot should not be higher than checkrange
	dbCheckPeriot int   // seconds
	dbQueryCheckRange int // seconds
	dbPath string
	logPath = "includes/logs.log"
	logginLevel uint8
	httpListenPort string
	y map[string]interface{}
	config map[string]string
)

func main(){
	DeclarePriorityQueue()
	CreateDbTables()
	go HttpWorker()
	go ConsumeFromDb()
	ConsumeFromQueue()
	fmt.Println("Program successfully started...")
	ch := make(chan bool)
	<-ch
}

// go calls init() func before everything.
func init(){
	CreateLogFile()

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
	rabbitMqConnection = config["rabbitMqConnection"]
	currentDriver = config["currentDriver"]
	httpListenPort = config["httpListenPort"]
	priorityQueueName = config["priorityQueueName"]
	lv, err := strconv.Atoi(config["loggingLevel"])
	if LoggingChecking(err,"Couldn't read loggingLevel in config.txt","",3){
		os.Exit(1)
	}
	logginLevel = uint8(lv)
	priorityRange, err = strconv.ParseInt(config["priorityRange"], 10, 64)
	if LoggingChecking(err,"Couldn't read priorityRange in config.txt","",3){
		os.Exit(2)
	}
	dbCheckPeriot,err = strconv.Atoi(config["dbCheckPeriot"])
	if LoggingChecking(err,"Couldn't read dbCheckPeriot in config.txt","",3) {
		os.Exit(3)
	}
	dbQueryCheckRange, err = strconv.Atoi(config["dbQueryCheckRange"])
	if LoggingChecking(err,"Couldn't read dbQueryCheckRange in config.txt","",3) {
		os.Exit(4)
	}
}

func ReadConfig(filename_fullpath string) map[string]string {
	prg := "ReadConfig()"
	file, err := os.Open(filename_fullpath)
	if LoggingChecking(err,fmt.Sprintf("%s: os.Open(): %s\n", prg, err)," Config file opened.",3){
		os.Exit(5)
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
			if LoggingChecking(err,fmt.Sprintf("%s: regexp.Compile(): error=%s",prg, err),"",2) {
			} else {
				confOption := re.FindStringSubmatch(line)[1]
				confValue := re.FindStringSubmatch(line)[2]
				options[confOption] = confValue
				LoggingChecking(nil,"",fmt.Sprintf("%s: Config Option = %-18s, Config Value = %-18s\n",prg,confOption,confValue),0)
			}
		}
	}
	//LoggingChecking(nil,"",fmt.Sprintf("%s: options[]: %+v\n", prg, options),0)
	if LoggingChecking(scanner.Err(),fmt.Sprintf("%s: scanner.Err(): %s\n", prg, scanner.Err()),"",3){
		os.Exit(6)
	}
	return options
}
