package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/transfashion/evoucher/libs"
	"github.com/transfashion/evoucher/libs/helper"
)

type ResultPageData struct {
	Nominal     string
	ExpiredDate string
	VoucherLink string
}

func (hdr *Handler) Result(w http.ResponseWriter, r *http.Request) {
	basehref := r.Header.Get("BASE_HREF")

	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "result"

	cdb := libs.CustomerDb
	vdb := libs.VoucherDb
	// qcs := libs.Qiscus

	linkreq := r.URL.Query().Get("reqid")
	ld, err := cdb.GetLinkRequestData(linkreq)
	if err != nil {
		FormError(w, r, err)
		return
	}

	// ambil data voucher
	voucher, err := vdb.GetVoucher(ld.VouId)
	if err != nil {
		log.Printf("voucher %s not found", ld.VouId)
	}

	vou_id_url, err := helper.Encrypt(voucher.Id)
	if err != nil {
		log.Println("gagal ekripsi voucher")
		FormError(w, r, fmt.Errorf("gagal medapapatkan kode enkripsi voucher"))
		return
	}

	voucherlik := fmt.Sprintf("%sview/%s", basehref, vou_id_url)

	//fmt.Println(voucher)
	data := &ResultPageData{
		Nominal:     humanize.Comma(int64(voucher.Value)),
		ExpiredDate: voucher.ExpiredDate.Format("2006-01-02"),
		VoucherLink: voucherlik,
	}

	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}
