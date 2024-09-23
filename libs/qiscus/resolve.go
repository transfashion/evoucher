package qiscus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (q *Qiscus) Resolve(room_id string) error {
	baseurl := q.Config.BaseUrl
	appcode := q.Config.AppCode
	secret := q.Config.Secret

	url := fmt.Sprintf("%s/api/v1/admin/service/mark_as_resolved", baseurl)
	method := "POST"

	notes := "resolved%%20by%%20server"
	msg := fmt.Sprintf("room_id=%s&notes=%s", room_id, notes)

	payload := strings.NewReader(msg)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Qiscus-App-Id", appcode)
	req.Header.Add("Qiscus-Secret-Key", secret)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer res.Body.Close()

	// body, err := io.ReadAll(res.Body)
	_, err = io.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
