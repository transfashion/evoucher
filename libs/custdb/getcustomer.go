package custdb

import (
	"database/sql"
)

func (c *CustomerDB) GetCustomer(phonenumber string) (bool, *Customer, error) {

	var customer Customer
	row := c.Connection.QueryRow("SELECT custwa_id, custwa_name, custwa_gender FROM mst_custwa WHERE custwa_id = ?", phonenumber)
	err := row.Scan(&customer.PhoneNumber, &customer.Name, &customer.Gender)
	if err == sql.ErrNoRows {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}
	return true, &customer, nil
}
