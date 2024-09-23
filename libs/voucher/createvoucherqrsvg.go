package voucher

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
	s.Rect(0, 0, image_width, image_height, "fill:white;stroke:none") // buat bacground
	qs.WriteQrSVG(s)                                                  // buat QR nya

	// buat brand logo
	if v.HeaderLogoData != nil {
		encodedData := base64.StdEncoding.EncodeToString(v.HeaderLogoData)
		s.Image(0, 0, qr_width, voulogo_height, "data:image/png;base64,"+encodedData)
	}

	// buat text kode voucher
	tx := int(image_width / 2)
	ty := voulogo_height + qr_height + 15
	s.Text(tx, ty, text, "black; font-size:20px; text-anchor:middle; font-family:monospace")

	// buat keterangan
	ty = ty + 25
	s.Text(tx, ty, "voucher potongan harga", "black;text-anchor:middle; font-size:14px; font-family:monospace")

	// buat nilai voucher
	vouchervalue := humanize.Comma(int64(v.Value))

	ty = ty + 40
	s.Text(tx, ty, vouchervalue, "black; text-anchor:middle; font-size:36px; font-weight:bold; font-family:monospace")

	// tampilkan expired date
	// voucher.ExpiredDate.Format("2006-01-02"),
	expdate := fmt.Sprintf("voucher valid s/d %s", v.ExpiredDate.Format("02 Jan 2006"))
	ty = ty + 23
	s.Text(tx, ty, expdate, "black; text-anchor:middle; font-size:12px; font-family:monospace")

	// buat nilai tnc
	ty = ty + 15
	s.Text(tx, ty, "** syarat dan kententuan berlaku **", "black; text-anchor:middle; font-size:9px; font-style:italic; font-family:monospace")

	s.End()

	return buffer.String(), nil
}
