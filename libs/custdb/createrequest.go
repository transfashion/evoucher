package custdb

import (
	"fmt"

	"github.com/transfashion/evoucher/libs/uniqid"
)

type RequestData struct {
	Customer   *Customer
	Ref        string
	RoomId     string
	VoubatchId string
	Intent     string
	Message    string
	JsonData   string
}

func (c *CustomerDB) CreateRequest(req *RequestData) (string, error) {
	uiq := uniqid.New(uniqid.Params{MoreEntropy: true})[:14]
	reqid := addParity(uiq)

	// simpan ke database
	query := `
		insert into mst_custwalinkreq
		(custwalinkreq_id, ref, intent, room_id, message, data, voubatch_id, custwa_id, _createby)
		values
		(?, ?, ?, ?, ?, ?, ?, ?, '5effbb0a0f7d1')
	`
	_, err := c.Connection.Exec(query, reqid, req.Ref, req.Intent, req.RoomId, req.Message, req.JsonData, req.VoubatchId, req.Customer.PhoneNumber)
	if err != nil {
		return "", err
	}

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
