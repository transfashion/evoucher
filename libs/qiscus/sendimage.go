package qiscus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (q *Qiscus) SendImage(room_id string, imagelink string, message string) (string, error) {
	baseurl := q.Config.BaseUrl
	appcode := q.Config.AppCode
	sender := q.Config.Sender
	url := fmt.Sprintf("%s/%s/bot", baseurl, appcode)
	method := "POST"

	data := fmt.Sprintf(`{
		"sender_email": "%s", 
		"message": "%s",
		"type": "file_attachment",
		"room_id": "%s",
		"payload": {
			"url": "%s",
			"caption": "%s"
		}
	}`, sender, message, room_id, imagelink, message)

	log.Println("Send image via qiscus", imagelink)
	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("QISCUS_SDK_SECRET", q.Config.Secret)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	log.Println("result:", string(body))
	return string(body), nil
}
