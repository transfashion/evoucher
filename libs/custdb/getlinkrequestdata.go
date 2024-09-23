package custdb

import (
	"database/sql"
	"log"

	"github.com/transfashion/evoucher/libs/helper"
)

func (db *CustomerDB) GetLinkRequestData(reqid string) (*RequestData, error) {

	query := `
		select 
		A.cust_id,
		B.custwa_id, A.cust_name, ifnull(A.gender_id, '-') as gender, 
		B.ref, B.intent, B.room_id, B.message, B.data, B.voubatch_id, B.vou_id
		from mst_cust A
		join mst_custwalinkreq B on A.cust_id = B.cust_id
		where B.custwalinkreq_id = ?
	
	`

	var d RequestData
	var c Customer

	var vou_id *string = new(string)

	log.Println("GetLinkRequestData", reqid)
	row := db.Connection.QueryRow(query, reqid)
	err := row.Scan(&c.Id, &c.PhoneNumber, &c.Name, &c.Gender, &d.Ref, &d.Intent, &d.RoomId, &d.Message, &d.JsonData, &d.VoubatchId, &vou_id)
	if err == sql.ErrNoRows {
		log.Println("data request tidak ditemukan")
		return nil, nil
	} else if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	d.VouId = helper.IsStringNil(vou_id, "")

	ld := &d
	ld.Customer = &c

	return ld, nil
}
