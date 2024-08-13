package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/transfashion/evoucher/libs"
)

func (hdr *Handler) VoucherQrJPG(w http.ResponseWriter, r *http.Request) {
	vou_id := chi.URLParam(r, "vouid")

	vdb := libs.VoucherDb
	voucher, err := vdb.GetVoucher(vou_id)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		panic(err)
	}

	if voucher == nil {
		fmt.Fprintln(w, "voucher not found")
		return
	}

	logofilepath := filepath.Join(hdr.Webservice.RootDir, "data", "images", "vlogo_transfashion.png")
	logodata, err := os.ReadFile(logofilepath)
	if err != nil {
		log.Fatal(err)
	}

	if logodata != nil {
		voucher.HeaderLogoData = logodata
	}

	jpgdata, err := voucher.CreateVoucherQrJPG()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		panic(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(jpgdata)

}
