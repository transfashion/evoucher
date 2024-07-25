package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
)

type UrlQueryVoucher struct {
	RoomId      string `json:"room_id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type FormPageData struct {
	PhoneNumber string
	Name        string
	Gender      string
	Code        string

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

	var urlq UrlQueryVoucher
	err = json.Unmarshal(query, &urlq)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	data := &FormPageData{
		PhoneNumber: urlq.PhoneNumber,
		Name:        urlq.Name,
	}

	if r.Method == "POST" {
		var invalid bool = false

		phone := r.FormValue("phone")
		name := r.FormValue("name")
		gender := r.FormValue("gender")
		code := r.FormValue("code")

		data.Name = name
		data.PhoneNumber = phone
		data.Gender = gender
		data.Code = code

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
			fmt.Println(basehref)
			nexturl := fmt.Sprintf("%sresult", basehref)
			fmt.Println(nexturl)
			http.Redirect(w, r, nexturl, http.StatusSeeOther)
		}

	}

	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}

// func (hdr *Handler) FormPost(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
// 	pv.PageName = "form"
// 	defaulthandlers.SimplePageHandler(pv, w, r)
// }
