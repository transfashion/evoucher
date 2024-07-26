package qiscus

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (q *Qiscus) SendMessage(room_id string, message string) error {
	baseurl := q.Config.BaseUrl
	appcode := q.Config.AppCode
	sender := q.Config.Sender
	url := fmt.Sprintf("%s/%s/bot", baseurl, appcode)
	method := "POST"

	data := fmt.Sprintf(`{
		"sender_email": "%s", 
		"type": "text",
		"room_id": "%s",
		"message": "%s"
	}`, sender, room_id, message)

	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("QISCUS_SDK_SECRET", q.Config.Secret)

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
