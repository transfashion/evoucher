package custdb

func (c *CustomerDB) CommitLinkRequest(linkreq string, res string) error {
	query := `
		update mst_custwalinkreq set result=?, status='OK' where custwalinkreq_id = ?
	`
	_, err := c.Connection.Exec(query, res, linkreq)
	if err != nil {
		return err
	}
	return nil
}
