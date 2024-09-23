package voucher

import (
	"log"

	"github.com/fgtago/fgweb/appsmodel"
	"github.com/mskrha/svg2png"
	"github.com/transfashion/evoucher/models"
)

func (v *Voucher) CreateVoucherQrPNG() ([]byte, error) {
	ws := appsmodel.GetWebservice()
	appcfg := ws.ApplicationConfig.(*models.ApplicationConfig)

	svgdata, err := v.CreateVoucherQrSvg()
	if err != nil {
		return nil, err
	}

	converter := svg2png.New()
	if appcfg.Binaries.Inkscape != "" {
		//Mac OS converter.SetBinary("/Applications/Inkscape.app/Contents/MacOS/inkscape")
		converter.SetBinary(appcfg.Binaries.Inkscape)
	}

	input := []byte(svgdata)
	output, err := converter.Convert(input)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return output, nil
}
