package custdb

import "log"

func (c *CustomerDB) UpdateLinkRequestCode(linkreq string, code string, slpart string) error {
	query := `update mst_custwalinkreq set code=?, slpart=? where custwalinkreq_id = ?`
	_, err := c.Connection.Exec(query, code, slpart, linkreq)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
