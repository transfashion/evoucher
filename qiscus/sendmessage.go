package qiscus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (q *Qiscus) SendMessage(chat *Chat, message string) error {
	baseurl := q.BaseUrl
	appcode := q.AppCode
	sender := q.Sender
	url := fmt.Sprintf("%s/%s/bot", baseurl, appcode)
	method := "POST"

	data := fmt.Sprintf(`{
		"sender_email": "%s", 
		"type": "text",
		"room_id": "%s",
		"message": "%s"
	}`, sender, chat.RoomId, message)

	fmt.Println("sending message")
	fmt.Println(url)
	fmt.Println(data)

	payload := strings.NewReader(data)

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
