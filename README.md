# MailWorker
This project is about taking mail informations over **HTTP** protocol in **json** format and send them to the target **email** addresses in order with their **_priority_** or **_time_** tag.

## How to use ?
At this current version, only thing you need to do is get in the project directory and give that command at the console <br />**go run \*.go**.
Then the program will be started to listening incoming **POST** requests.

Different types of configuration settings can be set under config.txt file. Current http listening port is 9000.
Sample json format is given at the mail_ornekleri.txt. In this format, it's possible to send multiple email contents within the one json packet.
To use sample POST request sender and see if it works, get in the _includes_ subfolder and call it with url arg
 <br />**Example use:**  `go run send_post_request.go http://localhost:9000/post`

#### **Notes:** 
- Make sure RabbitMQ server is up and config.txt file is under ./includes/ directory.
- Dont use whitespace while editing config file.  \(Correct: option=value  , Wrong: **_option = value_** \)
- Current json format sat for **Sendgrid** api. \(simple example is right below.\)
- Program is not sending mails to targets, just showing at the logs. If you want to send or add a new driver, check out **_StartConsume_** function from **queue_consumer.go** file and add whatever action you want.


##### Sendgrid mail sender sample:
````package main
  
  import (
  	"fmt"
  	"github.com/sendgrid/sendgrid-go"
  	"log"
  )
  
  func main() {
  	request := sendgrid.GetRequest("SENDGRID-API-KEY-HERE", "/v3/mail/send", "https://api.sendgrid.com")
  	request.Method = "POST"
  	request.Body = []byte(
  `{
    "personalizations": [
      {
        "to": [
          {
            "email": "example1@gmail.com"
          }
        ],
        "subject": "Sending with SendGrid is Fun"
      }
    ],
    "from": {
      "email": "cengizucgul@hotmail.com",
      "name":"Cengiz Üçgül"
    },
    "content": [
      {
        "type": "text/plain",
        "value": "and easy to do anywhere, even with Go"
      }
    ],
    "reply_to": {
      "name":"Cengiz",
      "email":"cengizucgul07@gmail.com"
    },
    "attachments": [
      {
        "content": "c2FtcGxlIHRleHQgc3RyaW5ncy4uLg==",  // Base64 format
        "filename": "some-attachment.txt",
        "type": "plain/text",
        "disposition": "attachment",
        "contentId": "mytext"
      }
    ]
  }`)
  	response, err := sendgrid.API(request)
  	if err != nil {
  		log.Println(err)
  	} else {
  		fmt.Println(response.StatusCode)
  		fmt.Println(response.Body)
  		fmt.Println(response.Headers)
  	}
  }
 ````