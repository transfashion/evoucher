package voucher

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/transfashion/evoucher/libs/qrcode"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"

	"github.com/dustin/go-humanize"
)

func (v *Voucher) CreateVoucherQrSvg() (string, error) {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)

	blocksize := 10
	text := v.Id

	s := svg.New(writer)

	qrCode, _ := qr.Encode(text, qr.M, qr.Auto)

	// Write QR code to SVG
	qs := qrcode.NewQrSVG(qrCode, blocksize)
	qr_width := qs.GetImageWidth()
	qr_height := qr_width

	voulogo_height := 60
	vounote_height := 120

	image_width := qr_width
	image_height := voulogo_height + qr_height + vounote_height

	// buat QR
	qs.SetStartPoint(0, voulogo_height)
	s.Start(image_width, image_height)
	qs.WriteQrSVG(s)

	// baca logo

	// buat brand logo
	if v.HeaderLogoData != nil {
		encodedData := base64.StdEncoding.EncodeToString(v.HeaderLogoData)
		s.Image(0, 0, qr_width, voulogo_height, "data:image/png;base64,"+encodedData)
	}

	// buat note voucher
	tx := int(image_width / 2)
	ty := voulogo_height + qr_height + 18
	s.Text(tx, ty, text, "black;text-anchor:middle; font-size:24px")

	// buat keterangan
	ty = ty + 26
	s.Text(tx, ty, "voucher potongan harga", "black;text-anchor:middle; font-size:16px")

	// buat nilai voucher
	vouchervalue := humanize.Comma(int64(v.Value))

	ty = ty + 40
	s.Text(tx, ty, vouchervalue, "black; text-anchor:middle; font-size:42px; font-weight:bold")

	// buat nilai tnc
	ty = ty + 26
	s.Text(tx, ty, "** syarat dan kententuan berlaku **", "black; text-anchor:middle; font-size:12px; font-style:italic")

	s.End()

	return buffer.String(), nil
}
