package voucher

import (
	"fmt"
	"strconv"
)

func (v *VoucherDB) VerifyCode(code string) (string, bool) {

	_, err := strconv.Atoi(code)
	if err != nil {
		return "", false
	}

	rc, parity, parityfactor := v.parseCode(code) // ambil kode parity

	urc := v.reverseString(rc)
	irc, _ := strconv.Atoi(urc)
	nrc := irc - parityfactor

	crc := fmt.Sprintf("%05d", nrc)
	orc := fmt.Sprintf("%05d", nrc-12300)

	p := v.getParity(crc)

	if p != parity {
		return "", false
	} else {
		return orc, true
	}
}

func (v *VoucherDB) parseCode(code string) (rc string, parity int, parityfactor int) {
	lastchar := string(code[len(code)-1])
	parity, _ = strconv.Atoi(lastchar)

	pf := fmt.Sprintf("%d0%d", parity, parity)
	parityfactor, _ = strconv.Atoi(pf)

	rc = string(code[0 : len(code)-1])

	return rc, parity, parityfactor
}

func (v *VoucherDB) reverseString(s string) string {
	var reversed string
	for i := len(s) - 1; i >= 0; i-- {
		reversed += string(s[i])
	}
	return reversed
}

func (v *VoucherDB) getParity(str string) int {
	t := 9
	for i := 0; i < len(str); i++ {
		c := string(str[i])
		p, _ := strconv.Atoi(c)

		if p == 0 {
			t = t * (i * i)
		} else {
			t = t * (p * p)
		}
	}

	value := fmt.Sprintf("%d", t)
	lastchar := string(value[len(value)-1])
	parity, _ := strconv.Atoi(lastchar)

	return parity
}
