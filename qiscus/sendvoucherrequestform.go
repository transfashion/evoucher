package qiscus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (q *Qiscus) SendVoucherRequestForm(chat *Chat) error {
	baseurl := q.BaseUrl
	appcode := q.AppCode
	url := fmt.Sprintf("%s/%s/bot", baseurl, appcode)
	method := "POST"
	sender := q.Sender

	fmt.Println("Sending Voucher Request Form")
	fmt.Println(url)
	fmt.Println(chat)

	payload := strings.NewReader(fmt.Sprintf(`{
	  "sender_email": "%s", 
	  "message": "Hi good morning",
	  "type": "buttons",
	  "room_id": "%s",
	  "payload": {
		  "text": "Untuk Request Voucher silakan isi form dengan klik tombol dibawah ini",
		  "buttons": [
			  {
				  "label": "Request Voucher",
				  "type": "postback",
				  "payload": {
					  "url": "http://somewhere.com/button1",
					  "method": "get",
					  "payload": null
				  }
			  }
		  ]
	  }
  	}`, sender, chat.RoomId))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("QISCUS_SDK_SECRET", q.Secret)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
