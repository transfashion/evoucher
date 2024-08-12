package voucher

import (
	"github.com/mskrha/svg2png"
)

func (v *Voucher) CreateVoucherQrPNG() ([]byte, error) {
	svgdata, err := v.CreateVoucherQrSvg()
	if err != nil {
		return nil, err
	}

	converter := svg2png.New()
	//converter.SetBinary("/Applications/Inkscape.app/Contents/MacOS/inkscape")

	input := []byte(svgdata)
	output, err := converter.Convert(input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
