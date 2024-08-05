package custdb

import (
	"fmt"

	"github.com/transfashion/evoucher/libs/uniqid"
)

type RequestData struct {
	Customer *Customer
	RoomId   string
	Ref      string
	Intent   string
}

func (c *CustomerDB) CreateRequest(data *RequestData) (string, error) {
	uiq := uniqid.New(uniqid.Params{MoreEntropy: true})[:14]
	reqid := addParity(uiq)

	// simpan ke database

	return reqid, nil
}

func addParity(reqid string) string {
	parityValue := calculateParity(reqid)
	return fmt.Sprintf("%s%s", reqid, parityValue)
}

func calculateParity(data string) string {
	parity := byte(0)
	for _, char := range data {
		parity ^= byte(char)
	}
	parityHex := fmt.Sprintf("%02x", parity)
	return parityHex
}

/* untuk cek parity

str := "66b0886b70521700"
mv := str[:14]
pv := str[len(str)-2:]
fmt.Println(mv, pv)

// cek parity
parity := CalculateParity(mv)
fmt.Println(mv, parity)


*/
