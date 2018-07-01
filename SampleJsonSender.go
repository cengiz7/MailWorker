package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
)

func main() {
	url := "http://localhost:9000/post"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{
  "mails": [
    {
      "priority":"1",
      "content":{
        "headers": {
          "Date": "Mon, 16 Jan 2012 17:00:01 +0000",
          "From": "Message Sender <sender@example.com>",
          "To": ["mail1@google.com","mail2@google.com"],
          "Subject": "Test Subject",
          "Mime-Version": "1.0",
          "Content-Type": "multipart/alternative; boundary=------------090409040602000601080801",
          "User-Agent": "Postbox 3.0.2 (Macintosh/20111203)"},
        "body":{
          "text/html": "<html><head>\n<meta http-equiv=\"content-type\" content=\"text/html; charset=ISO-8859-1\"></head><body\n bgcolor=\"#FFFFFF\" text=\"#000000\">\nTest with <span style=\"font-weight: bold;\">HTML</span>.<br>\n</body>\n</html>"},
        "attachments": [
          {
            "file_name": "file1.txt",
            "content-type": "text/plain",
            "size": 8,
            "url": "http://example.com/file1.txt",
            "disposition": "attachment"
          },
          {
            "file_name": "file.txt",
            "content_type": "text/plain",
            "size": 8,
            "content": "dGVzdGZpbGU=",
            "disposition": "attachment"
          }
        ]
      }
    }
  ]
}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}