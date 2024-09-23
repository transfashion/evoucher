package custdb

import (
	"database/sql"
	"log"
)

func (c *CustomerDB) GetCustomer(phonenumber string) (bool, *Customer, error) {

	var customer Customer
	// row := c.Connection.QueryRow("SELECT custwa_id, custwa_name, custwa_gender FROM mst_custwa WHERE custwa_id = ?", phonenumber)

	query := `
		select 
			A.cust_id as custid,
			B.custaccess_code as phonenumber,
			A.cust_name as custname,
			ifnull(A.gender_id, '-') as gender
		from
		mst_cust A inner join mst_custaccess B on B.cust_id=A.cust_id 
		where 
				B.custaccesstype_id = 'WA'
			and B.custaccess_code  = ?
	`

	row := c.Connection.QueryRow(query, phonenumber)
	err := row.Scan(&customer.Id, &customer.PhoneNumber, &customer.Name, &customer.Gender)
	if err == sql.ErrNoRows {
		log.Println(err.Error())
		return false, nil, nil
	} else if err != nil {
		log.Println(err.Error())
		return false, nil, err
	}
	return true, &customer, nil
}
