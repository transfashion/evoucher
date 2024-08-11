package custdb

import "database/sql"

func (c *CustomerDB) GetPendingRequest(cust *Customer, ref string) (string, error) {
	// cari request dari customer ini dengan ref, yang statusnya masih null
	query := `select custwalinkreq_id from mst_custwalinkreq where custwa_id = ? and ref = ? and status is null`

	var reqid string
	row := c.Connection.QueryRow(query, cust.PhoneNumber, ref)
	err := row.Scan(&reqid)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return reqid, nil
}
