package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/transfashion/evoucher/libs"
)

type FormPageData struct {
	RequestId   string
	PhoneNumber string
	Name        string
	Gender      string
	Code        string
	RoomId      string

	GenderInvalid      bool
	CodeInvalid        bool
	CodeInvalidMessage string
	SLPart             string
}

func (hdr *Handler) Form(w http.ResponseWriter, r *http.Request) {
	basehref := r.Header.Get("BASE_HREF")

	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "form"

	cdb := libs.CustomerDb
	vdb := libs.VoucherDb
	qcs := libs.Qiscus

	linkreq := r.URL.Query().Get("req")
	ld, err := cdb.GetLinkRequestData(linkreq)
	if err != nil {
		FormError(w, r, err)
		return
	}

	if ld == nil {
		FormError(w, r, fmt.Errorf("link tidak ditemukan"))
		return
	}

	data := &FormPageData{
		RequestId:   linkreq,
		RoomId:      ld.RoomId,
		PhoneNumber: ld.Customer.PhoneNumber,
		Name:        ld.Customer.Name,
	}

	if r.Method == "POST" {
		log.Println("posting data")
		var invalid bool = false

		phone := r.FormValue("phone")
		name := r.FormValue("name")
		gender := r.FormValue("gender")
		code := r.FormValue("code")
		room_id := r.FormValue("room_id")
		request_id := r.FormValue("request_id")

		data.Name = name
		data.PhoneNumber = phone
		data.Gender = gender
		data.Code = code
		data.RoomId = room_id
		data.RequestId = request_id

		// cek gender sudah diisi apa belum
		if gender == "" {
			invalid = invalid || true
			data.GenderInvalid = true
		}

		// cek apakah code sudah diisi
		if data.Code == "" {
			invalid = invalid || true
			data.CodeInvalid = true
			data.CodeInvalidMessage = "Code harus diisi. Silakan minta kode aktifasi voucher ke kasir"
		} else {
			// cek apakah code yang diisikan sudah benar
			/* contoh kode valid
			018819
			805816
			400811
			803814
			014815
			*/

			log.Println("verifying code", data.Code)
			slpart, isvalidcode := vdb.VerifyCode(data.Code)
			if !isvalidcode {
				invalid = invalid || true
				data.CodeInvalid = true
				data.CodeInvalidMessage = "Code yang diisikan salah"
			}
			data.SLPart = slpart
		}

		if !invalid {
			log.Println("Processing Data")
			// data yang diisikan sudah benar

			// save data customer
			err := cdb.UpdateCustomer(phone, name, gender)
			if err != nil {
				FormError(w, r, err)
				return
			}

			// update linkrequest
			log.Println("update linkrequest code", linkreq, code)
			err = cdb.UpdateLinkRequestCode(linkreq, data.Code, data.SLPart)
			if err != nil {
				FormError(w, r, err)
				return
			}

			// cek apakah voucher untuk code dan phone number udah di create
			// voucher_id, exists, err := vdb.GetVoucher(ld.VoubatchId, ld.Customer.PhoneNumber, code)
			log.Println("Get existing user voucher code", ld.VoubatchId, code, "for", ld.Customer.PhoneNumber, code)
			voucher, err := vdb.GetVoucherByPhoneNumber(ld.VoubatchId, ld.Customer.PhoneNumber, code)
			if err != nil {
				FormError(w, r, err)
				return
			}

			if voucher != nil {
				// voucher telah dibuat, redirect ke halaman preview voucher
				log.Println("voucher already created", voucher.Id)

				// send information to whatsapp
				log.Println("sending message via qiscus to", data.PhoneNumber, data.RoomId)
				//voucherlik := fmt.Sprintf("%s%s/voucherqr.svg", basehref, voucher.Id)
				voucherlik := fmt.Sprintf("%sview/%s", basehref, voucher.Id)
				_, errqcs := qcs.SendMessage(data.RoomId, fmt.Sprintf("Anda telah mempunyai voucher ini dari request sebelumnya. Silakan klik link %s untuk melihat voucher anda", voucherlik))
				if errqcs != nil {
					FormError(w, r, errqcs)
					return
				}

				// redirect ke halaman preview voucher
				nexturl := fmt.Sprintf("%vsent", basehref)
				http.Redirect(w, r, nexturl, http.StatusSeeOther)
				return
			}

			// customer belum mempunyai voucher untuk batch dan kode ini
			// buat voucher baru
			voucher, err = vdb.CreateNewVoucher(ld.VoubatchId, ld.Customer.PhoneNumber, ld.Customer.Name, code)
			if err != nil {
				FormError(w, r, err)
				return
			}

			log.Println("new voucher issued", voucher.Id)

			// update ke ke mst_custwalinkreq untuk kode voucher dan value ke field result
			log.Println("update linkrequest voucher", linkreq, code)
			err = cdb.UpdateLinkRequestVoucher(linkreq, voucher.Id)
			if err != nil {
				FormError(w, r, err)
				return
			}

			// kirimkan informasi ke whatsapp untuk kode voucher

			// send image voucher
			imglink := "https://evoucher.transfashionindonesia.com/testqr.svg"
			//imglink := fmt.Sprintf("%s%s/voucherqr.svg", basehref, voucher.Id)
			log.Println("sending image voucher via qiscus to", data.PhoneNumber, data.RoomId, imglink)
			res, err := qcs.SendImage(data.RoomId, imglink, "Tunjukkan voucher ini saat bertransaksi untuk mendapatkan potongan harga senilai voucher. (Syarat dan ketentuan berlaku)")
			if err != nil {
				FormError(w, r, err)
				return
			}
			log.Println(res)

			// send message
			// temphref := basehref
			// basehref = "https://evoucher.transfashionindonesia.com/"
			log.Println("sending message via qiscus to", data.PhoneNumber, data.RoomId)
			voucherlik := fmt.Sprintf("%sview/%s", basehref, voucher.Id)
			tmp := "Hai kak %s, selamat anda mendapatkan voucher potongan harga senilai %s. Untuk melihat dan menggunakan voucher ini, bisa juga dengan klik link %s. Terimakasih."
			vouvalue := humanize.Comma(int64(voucher.Value))
			msg := fmt.Sprintf(tmp, data.Name, vouvalue, voucherlik)
			res, err = qcs.SendMessage(data.RoomId, msg)
			if err != nil {
				FormError(w, r, err)
				return
			}
			// basehref = temphref

			// commit linkrequest
			err = cdb.CommitLinkRequest(linkreq, res)
			if err != nil {
				FormError(w, r, err)
				return
			}

			// 	// redirect
			nexturl := fmt.Sprintf("%sresult?reqid=%s", basehref, linkreq)
			http.Redirect(w, r, nexturl, http.StatusSeeOther)
			return
		} else {
			log.Println("invalid data")
		}

	}

	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}

func FormError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
