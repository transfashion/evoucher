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
	TNC              []string
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

	// get voucher TNC
	query := "select voutnc_descr from mst_voutnc  where voubatch_id = ? order by voutnc_order"
	rows, err := vdb.Connection.Query(query, voucher.BatchId)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		panic(err)
	}
	defer rows.Close()

	tncs := make([]string, 0)
	for rows.Next() {
		var tnc string
		err := rows.Scan(&tnc)
		if err != nil {
			fmt.Fprintln(w, err.Error())
			panic(err)
		}
		tncs = append(tncs, tnc)
	}

	data := &ViewPageData{
		VoucherId:        vou_id,
		VoucherQrSvgLink: voucherlik,
		VoucherTNC:       "tnc dari voucher",
		TNC:              tncs,
	}

	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "view"
	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}
