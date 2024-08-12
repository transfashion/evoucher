package custdb

import (
	"database/sql"
	"log"
)

func (db *CustomerDB) GetLinkRequestData(reqid string) (*RequestData, error) {

	query := `
		select 
		A.custwa_id, A.custwa_name, A.custwa_gender, 
		B.ref, B.intent, B.room_id, B.message, B.data, B.voubatch_id, B.vou_id
		from mst_custwa A
		join mst_custwalinkreq B on A.custwa_id = B.custwa_id
		where B.custwalinkreq_id = ?
	
	`

	var d RequestData
	var c Customer

	log.Println("GetLinkRequestData", reqid)
	row := db.Connection.QueryRow(query, reqid)
	err := row.Scan(&c.PhoneNumber, &c.Name, &c.Gender, &d.Ref, &d.Intent, &d.RoomId, &d.Message, &d.JsonData, &d.VoubatchId, &d.VouId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	ld := &d
	ld.Customer = &c

	return ld, nil
}
