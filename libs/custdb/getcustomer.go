package custdb

func (c *CustomerDB) GetCustomer(phonenumber string) (bool, *Customer, error) {

	var customer Customer
	err := c.Connection.QueryRow("SELECT phonenumber, name FROM customers WHERE phonenumber = ?", phonenumber).Scan(&customer.PhoneNumber, &customer.Name)
	if err != nil {
		return false, nil, err
	}

	return false, nil, nil
}
