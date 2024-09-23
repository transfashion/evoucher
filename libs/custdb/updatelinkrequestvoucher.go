package custdb

import "log"

func (c *CustomerDB) UpdateLinkRequestVoucher(linkreq string, vou_id string) error {
	query := `
		update mst_custwalinkreq set vou_id=? where custwalinkreq_id = ?
	`
	_, err := c.Connection.Exec(query, vou_id, linkreq)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
