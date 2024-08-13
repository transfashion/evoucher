package voucher

import (
	"bytes"
	"image/jpeg"
	"image/png"

	"github.com/mskrha/svg2png"
)

func (v *Voucher) CreateVoucherQrJPG() ([]byte, error) {
	svgdata, err := v.CreateVoucherQrSvg()
	if err != nil {
		return nil, err
	}

	converter := svg2png.New()
	//converter.SetBinary("/Applications/Inkscape.app/Contents/MacOS/inkscape")

	input := []byte(svgdata)
	pngdata, err := converter.Convert(input)
	if err != nil {
		return nil, err
	}

	// convert png data to jpg
	output, err := JpegToPng(pngdata)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func JpegToPng(pngdata []byte) ([]byte, error) {

	// Decode file PNG
	img, err := png.Decode(bytes.NewReader(pngdata))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
