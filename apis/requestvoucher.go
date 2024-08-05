package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/transfashion/evoucher/libs"
	"github.com/transfashion/evoucher/libs/custdb"
	"github.com/transfashion/evoucher/libs/uniqid"
	"github.com/transfashion/evoucher/models"
)

type RequestVoucherPayload struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	RoomId      string `json:"room_id"`
	FromName    string `json:"from_name"`
	Intent      string `json:"intent"`
}

func (api *Api) RequestVoucher(w http.ResponseWriter, r *http.Request) {
	payload := RequestVoucherPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Intent == "#VoucherPromoTransfashion" {
		ws := api.Webservice
		appconf := ws.ApplicationConfig.(*models.ApplicationConfig)

		cdb := libs.CustomerDb
		qcs := libs.Qiscus
		vdb := libs.VoucherDb

		var exist bool
		var cust *custdb.Customer
		exist, cust, err = cdb.GetCustomer(payload.PhoneNumber)
		if err != nil {
			fmt.Println(err)
		}

		fromname := strings.ToValidUTF8(payload.FromName, "")
		if !exist {
			cdb.CreateNew(payload.PhoneNumber, fromname)
		}

		/* parse message */
		vdb.ParseMessage(payload.Message)

		query := models.FormUrlQuery{
			RequestId: uniqid.New(uniqid.Params{}),
			RoomId:    payload.RoomId,
			Number:    payload.PhoneNumber,
			Name:      payload.FromName,
			Batch:     "evoucher",
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

		msg := fmt.Sprintf("Halo %s, kamu mendapatkan voucher promo dari *Trans Fashion Indonesia*. Untuk aktifasi voucher, silahkan klik link berikut: %s/form?q=%s", payload.FromName, appconf.Evoucher.Url, q)
		err = qcs.SendMessage(payload.RoomId, msg)
		if err != nil {
			fmt.Println(err)
		}
		// data := ""

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\": true}"))
}
