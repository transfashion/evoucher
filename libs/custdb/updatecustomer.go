package custdb

import "log"

func (c *CustomerDB) UpdateCustomer(cust_id string, name string, gender string) error {
	query := `update mst_cust set cust_name=?, gender_id=? where cust_id=?`
	_, err := c.Connection.Exec(query, name, gender, cust_id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
