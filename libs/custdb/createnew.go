package custdb

import "fmt"

func (c *CustomerDB) CreateNew(phonenumber string, name string) (*Customer, error) {
	_, err := c.Connection.Exec("INSERT INTO mst_custwa (custwa_id, custwa_name, custwa_gender, _createby) VALUES (?, ?, ?, '5effbb0a0f7d1')", phonenumber, name, "-")
	if err != nil {
		return nil, err
	}

	exist, cust, err := c.GetCustomer(phonenumber)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("customer not inserted")
	}

	return cust, nil
}
