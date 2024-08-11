package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/transfashion/evoucher/libs"
	"github.com/transfashion/evoucher/libs/custdb"
	"github.com/transfashion/evoucher/models"
)

type RequestVoucherPayload struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	RoomId      string `json:"room_id"`
	FromName    string `json:"from_name"`
	Intent      string `json:"intent"`
}

const _VOUCHER_PROMO_INTENT_ = "#VoucherPromoTransfashion"

func (api *Api) RequestVoucher(w http.ResponseWriter, r *http.Request) {
	payload := RequestVoucherPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Voucher Request hit")
	if payload.Intent == _VOUCHER_PROMO_INTENT_ {
		ws := api.Webservice
		appconf := ws.ApplicationConfig.(*models.ApplicationConfig)

		cdb := libs.CustomerDb
		qcs := libs.Qiscus
		vdb := libs.VoucherDb

		var exist bool
		var cust *custdb.Customer
		log.Println("get customer with phone number", payload.PhoneNumber)
		exist, cust, err = cdb.GetCustomer(payload.PhoneNumber)
		if err != nil {
			RequestVoucherError(w, r, payload, err)
			return
		}

		fromname := strings.ToValidUTF8(payload.FromName, "")
		if !exist {
			log.Println("Customer", payload.PhoneNumber, "not found")
			log.Println("create new customer with phone number", payload.PhoneNumber)
			cust, err = cdb.CreateNew(payload.PhoneNumber, fromname)
			if err != nil {
				RequestVoucherError(w, r, payload, err)
				return
			}

		} else {
			log.Println("Customer found: ", cust.Name)
		}

		/* parse message */
		log.Println("Parse message")
		msgi := vdb.ParseMessage(payload.Message)
		if msgi == nil {
			RequestVoucherError(w, r, payload, fmt.Errorf("format request voucher tidak sesuai"))
			return
		}

		if msgi.VoubatchId == "" {
			RequestVoucherError(w, r, payload, fmt.Errorf("voucher batch code not found"))
			return
		}

		var reqid string

		/* get pending request */
		reqid, err = cdb.GetPendingRequest(cust, msgi.VoubatchId)
		if err != nil {
			RequestVoucherError(w, r, payload, err)
			return
		}

		if reqid == "" {
			/* request baru, masukkan request voucher ke database */
			log.Println("Create request voucher", msgi.VoubatchId, cust.PhoneNumber)
			reqid, err = cdb.CreateRequest(&custdb.RequestData{
				Customer:   cust,
				RoomId:     payload.RoomId,
				Intent:     payload.Intent,
				Ref:        msgi.VoubatchId,
				VoubatchId: msgi.VoubatchId,
				Message:    payload.Message,
				JsonData:   "{}",
			})
			if err != nil {
				RequestVoucherError(w, r, payload, err)
				return
			}
		}

		// cek apakah msgi.VoubatchId masih berlaku
		voubatch, err := vdb.GetVoucherBatch(msgi.VoubatchId)
		if err != nil {
			RequestVoucherError(w, r, payload, err)
			return
		}

		// verifikasi voucherbatch apakah valid:
		// - apakah periode generate masih berlaku ?
		// - dll
		_, err = vdb.VerifyBatch(voubatch)
		if err != nil {
			RequestVoucherError(w, r, payload, err)
			log.Println("sending message via qiscus")
			_, errqcs := qcs.SendMessage(payload.RoomId, err.Error())
			if errqcs != nil {
				RequestVoucherError(w, r, payload, errqcs)
				return
			}
			log.Println("message sent")
			return
		}

		/* kirimkan link form registrasi ke customer via qiscus */
		link := fmt.Sprintf("%s/form?req=%s", appconf.Evoucher.Url, reqid)
		msg := fmt.Sprintf("Halo %s, kamu mendapatkan voucher promo dari *Trans Fashion Indonesia*. Untuk aktifasi voucher, silahkan klik link berikut: %s", payload.FromName, link)
		_, err = qcs.SendMessage(payload.RoomId, msg)
		if err != nil {
			fmt.Println(err)
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\": true}"))
}

func RequestVoucherError(w http.ResponseWriter, r *http.Request, payload RequestVoucherPayload, err error) {
	log.Println(err.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"success\": false, \"error\": \"" + err.Error() + "\"}"))
}
