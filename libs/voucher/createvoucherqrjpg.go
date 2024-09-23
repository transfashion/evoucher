package voucher

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/mskrha/svg2png"
	"github.com/transfashion/evoucher/models"
)

func (v *Voucher) CreateVoucherQrJPG() ([]byte, error) {
	ws := appsmodel.GetWebservice()
	appcfg := ws.ApplicationConfig.(*models.ApplicationConfig)

	svgdata, err := v.CreateVoucherQrSvg()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	converter := svg2png.New()
	if appcfg.Binaries.Inkscape != "" {
		//Mac OS converter.SetBinary("/Applications/Inkscape.app/Contents/MacOS/inkscape")
		converter.SetBinary(appcfg.Binaries.Inkscape)
	}

	input := []byte(svgdata)
	pngdata, err := converter.Convert(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// convert png data to jpg
	output, err := JpegToPng(pngdata)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return output, nil
}

func JpegToPng(pngdata []byte) ([]byte, error) {

	// Decode file PNG
	img, err := png.Decode(bytes.NewReader(pngdata))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}
