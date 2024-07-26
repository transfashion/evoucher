package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/transfashion/evoucher/models"
)

type FormPageData struct {
	RequestId   string
	PhoneNumber string
	Name        string
	Gender      string
	Code        string
	RoomId      string

	GenderInvalid bool
	CodeInvalid   bool
}

func (hdr *Handler) Form(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "form"

	q := r.URL.Query().Get("q")
	query, err := base64.StdEncoding.DecodeString(q)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	var urlq models.FormUrlQuery
	err = json.Unmarshal(query, &urlq)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// cek apakah kode request ini telah dipakai
	// libs.VoucherDb.CreateVoucherRequest()

	data := &FormPageData{
		RequestId:   urlq.RequestId,
		RoomId:      urlq.RoomId,
		PhoneNumber: urlq.Number,
		Name:        urlq.Name,
	}

	if r.Method == "POST" {
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

		if code == "" {
			invalid = invalid || true
			data.CodeInvalid = true
		}

		if gender == "" {
			invalid = invalid || true
			data.GenderInvalid = true
		}

		if !invalid {
			// 	// save data

			// 	// redirect
			basehref := r.Header.Get("BASE_HREF")
			nexturl := fmt.Sprintf("%sresult", basehref)
			http.Redirect(w, r, nexturl, http.StatusSeeOther)
		}

	}

	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}
