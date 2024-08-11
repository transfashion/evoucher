package handlers

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/transfashion/evoucher/libs"

	"github.com/foolin/go-resvg"
)

func (hdr *Handler) VoucherQrPNG(w http.ResponseWriter, r *http.Request) {
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

	svgdata, err := voucher.CreateVoucherQrSvg()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		panic(err)
	}

	// myReader := strings.NewReader(svgdata)

	// convert SVG data to PNG

	// installasi module yang diperlukan di OS
	// sudo apt-get install libcairo2-dev
	// brew install cairo pkg-config
	img, err := resvg.RenderImageFromString(svgdata, nil)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	var buffer bytes.Buffer
	writer := io.Writer(&buffer)

	err = png.Encode(writer, img)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(buffer.Bytes())

}
