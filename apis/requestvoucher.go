package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/transfashion/evoucher/custdb"
	"github.com/transfashion/evoucher/models"
	"github.com/transfashion/evoucher/qiscus"
)

type RequestVoucherPayload struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	RoomId      string `json:"room_id"`
	FromName    string `json:"from_name"`
	Intent      string `json:"intent"`
}

func (api *Api) RequestVoucher(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

	payload := RequestVoucherPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat := &qiscus.Chat{
		RoomId: string(payload.RoomId),
		Number: string(payload.PhoneNumber),
		Name:   string(payload.FromName),
	}

	if payload.Intent == "#RequestVoucher" {
		ws := api.Webservice
		appconf := ws.ApplicationConfig.(*models.ApplicationConfig)
		appdata := ws.ApplicationData.(*models.ApplicationData)

		cdb := appdata.CustomerDb
		qcs := appdata.Qiscus

		var exist bool
		var cust *custdb.Customer
		exist, cust, err = cdb.GetCustomer(payload.PhoneNumber)
		if err != nil {
			fmt.Println(err)
		}

		query := models.FormUrlQuery{
			RoomId: payload.RoomId,
			Number: payload.PhoneNumber,
			Name:   payload.FromName,
		}

		if !exist {
			cdb.CreateNew(payload.PhoneNumber, payload.FromName)
		} else {
			query.Name = cust.Name
		}

		b, err := json.Marshal(query)
		if err != nil {
			fmt.Println(err)
			return
		}

		q := base64.StdEncoding.EncodeToString(b)
		link := fmt.Sprintf("%s/form?q=%s", appconf.Evoucher.Url, q)
		fmt.Println(link)

		msg := fmt.Sprintf("Silakan klik link dibawah untuk mendapatkan vouchermu %s/form?q=%s", appconf.Evoucher.Url, q)
		err = qcs.SendMessage(chat, msg)
		if err != nil {
			fmt.Println(err)
		}
		// data := ""
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// json.NewEncoder(w).Encode(data)
	}
}
