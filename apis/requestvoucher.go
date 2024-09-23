package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/transfashion/evoucher/libs"
	"github.com/transfashion/evoucher/libs/custdb"
	"github.com/transfashion/evoucher/libs/helper"
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
			log.Println(err.Error())
			RequestVoucherError(w, r, payload, err)
			return
		}

		fromname := strings.ToValidUTF8(payload.FromName, "")
		if !exist {
			log.Println("Customer", payload.PhoneNumber, "not found")
			log.Println("create new customer with phone number", payload.PhoneNumber)
			cust, err = cdb.CreateNew(payload.PhoneNumber, fromname)
			if err != nil {
				err := fmt.Errorf("failed to create new customer")
				log.Println(err.Error())
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
			err := fmt.Errorf("format request voucher tidak sesuai")
			log.Println(err.Error())
			RequestVoucherError(w, r, payload, err)
			return
		}

		if msgi.VoubatchId == "" {
			err := fmt.Errorf("voucher batch code not found")
			log.Println(err.Error())
			RequestVoucherError(w, r, payload, err)
			return
		}

		var reqid string

		/* get pending request */
		reqid, err = cdb.GetPendingRequest(cust, msgi.VoubatchId)
		if err != nil {
			log.Println(err.Error())
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
				log.Println(err.Error())
				RequestVoucherError(w, r, payload, err)
				return
			}
		}

		// ambil data voucherbatch
		voubatch, err := vdb.GetVoucherBatch(msgi.VoubatchId)
		if err != nil {
			log.Println(err.Error())
			RequestVoucherError(w, r, payload, err)
			return
		}

		// verifikasi voucherbatch apakah sudah active dan valid:
		// - apakah periode generate masih berlaku ?
		// - dll
		log.Println("verifying voucher batch", voubatch.Id)
		_, err = vdb.VerifyBatchRequest(voubatch)
		if err != nil {
			log.Println(err.Error())
			RequestVoucherError(w, r, payload, err)

			log.Println("sending informationmessage via qiscus")
			_, errqcs := qcs.SendMessage(payload.RoomId, err.Error())
			if errqcs != nil {
				log.Println("gagal kirim pesan via qiscus")
				RequestVoucherError(w, r, payload, errqcs)
				return
			}
			log.Println("message sent")
			return
		}

		// cek apakah customer sudah punya voucher untuk batch ini
		// jika sudah, kirimkan voucher link yang belum dipakai oleh customer via qiscus
		// pesan paling bawah: untuk aktivasi voucher lain, silakan klik link aktivasi

		// apabila bisa multiple voucher, ambil daftar voucher yang belum dipakai

		// jika singgle voucher, ambil voucher yang telah digenerate,
		// infokan apabila customer telah mempunyai voucher untuk batch ini,
		// dan infokan juga apabila sudah punya dan sudah dipakai

		var requestNewVoucher bool = false
		var replyMessage string
		if voubatch.IsMultiple {
			// multiple voucher
			vouchers, err := vdb.GetAllUnusedVoucherByPhoneNumber(msgi.VoubatchId, cust.PhoneNumber)
			if err != nil {
				log.Println("gagal get all unused voucher")
				RequestVoucherError(w, r, payload, err)
				return
			}

			if len(vouchers) > 0 {
				// Encrypt data voucher sebelum dikirim via URL
				voucher_id := vouchers[0].Id
				vou_id_url, err := helper.Encrypt(voucher_id)
				if err != nil {
					log.Println(err.Error())
					RequestVoucherError(w, r, payload, err)
					return
				}

				// customer sudah mempunyai voucher untuk batch ini
				// kirimkan link, untuk satu voucher ini, dan link aktivasi untuk menambah voucher berdasarkan code yang berbeda
				voucherlik := fmt.Sprintf("%s/view/%s", appconf.Evoucher.Url, vou_id_url)
				link := fmt.Sprintf("%s/form?req=%s", appconf.Evoucher.Url, reqid)
				replyMessage = fmt.Sprintln("Anda mempunyai voucher untuk batch ini.\nSilakan klik link berikut untuk melihat voucher:\n", voucherlik, "\n\nUntuk aktivasi voucher lain, silakan klik link berikut:\n", link)
			} else {
				// kirimkan link aktifasi voucher
				requestNewVoucher = true
			}

		} else {
			// single voucher
			voucher, err := vdb.GetSingleVoucherByPhoneNumber(msgi.VoubatchId, cust.PhoneNumber)
			if err != nil {
				log.Println("gagal get single voucher")
				RequestVoucherError(w, r, payload, err)
				return
			}

			if voucher != nil {
				// Encrypt data voucher sebelum dikirim via URL
				vou_id_url, err := helper.Encrypt(voucher.Id)
				if err != nil {
					log.Println(err.Error())
					RequestVoucherError(w, r, payload, err)
					return
				}

				// customer sudah mempunyai voucher untuk batch ini
				// kirimkan link, untuk voucher ini, berikan info voucher sudah dipakai atau belum
				if voucher.IsUse {
					replyMessage = fmt.Sprintf("Anda sudah mempunyai voucher untuk batch ini, dan sudah dipakai pada tanggal %s.\nTerimakasih", voucher.UsedDate)
				} else {
					voucherlik := fmt.Sprintf("%s/view/%s", appconf.Evoucher.Url, vou_id_url)
					replyMessage = fmt.Sprintf("Anda sudah mempunyai voucher untuk batch ini.\n\nSilahkan klik link berikut untuk melihat voucher ini:\n%s", voucherlik)
				}
			} else {
				// kirimkan link aktifasi voucher
				requestNewVoucher = true
			}

		}

		if requestNewVoucher {
			/* kirimkan link pembuatan voucher baru form registrasi ke customer via qiscus */
			var greeting string
			link := fmt.Sprintf("%s/form?req=%s", appconf.Evoucher.Url, reqid)
			if *voubatch.Greeting != "" && voubatch.Greeting != nil {
				greeting = fmt.Sprintf("Hai %s\n\n%s\n%s", payload.FromName, *voubatch.Greeting, link)
				// greeting = "test"
			} else {
				greeting = fmt.Sprintf("Hai %s\n\n, kamu mendapatkan voucher promo dari *Trans Fashion Indonesia*.\n\nUntuk aktifasi voucher, silahkan klik link berikut:\n%s", payload.FromName, link)
			}
			replyMessage = greeting
		}

		/* kirimkan message via qiscus */
		log.Println("sending message via qiscus")
		msg := strings.ReplaceAll(replyMessage, "\n", "\\n")
		_, err = qcs.SendMessage(payload.RoomId, msg)
		if err != nil {
			log.Println("gagal kirim pesan via qiscus")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"success\": true}"))

	} else {
		// intent tidak sesuai, pesan tidak diproses
		log.Printf("intent tidak sesuai (%s)\n", payload.Intent)
		res := fmt.Sprintf("{\"success\": false, \"message\": \"Message tidak di process, intent tidak sesuai (%s)\"}", payload.Intent)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}

}

func RequestVoucherError(w http.ResponseWriter, r *http.Request, payload RequestVoucherPayload, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\"success\": false, \"error\": \"" + err.Error() + "\"}"))
}
