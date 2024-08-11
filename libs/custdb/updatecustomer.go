package custdb

func (c *CustomerDB) UpdateCustomer(phonenumber string, name string, gender string) error {
	query := `update mst_custwa set custwa_name = ?, custwa_gender = ? where custwa_id = ?`
	_, err := c.Connection.Exec(query, name, gender, phonenumber)
	if err != nil {
		return err
	}
	return nil
}
