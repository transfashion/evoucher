package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/transfashion/evoucher/libs"

	"github.com/go-chi/chi/v5"
)

type ViewPageData struct {
	VoucherId        string
	VoucherQrSvgLink string
	VoucherTNC       string
}

func (hdr *Handler) View(w http.ResponseWriter, r *http.Request) {
	basehref := r.Header.Get("BASE_HREF")
	vdb := libs.VoucherDb

	vou_id := chi.URLParam(r, "vouid")

	// ambil data voucher
	voucher, err := vdb.GetVoucher(vou_id)
	if err != nil {
		log.Printf("voucher %s not found", vou_id)
	}

	if voucher == nil {
		nexturl := fmt.Sprintf("%snotfound", basehref)
		http.Redirect(w, r, nexturl, http.StatusSeeOther)
		return
	}

	voucherlik := fmt.Sprintf("%s%s/voucherqr.svg", basehref, vou_id)

	data := &ViewPageData{
		VoucherId:        vou_id,
		VoucherQrSvgLink: voucherlik,
		VoucherTNC:       "tnc dari voucher",
	}

	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "view"
	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}
