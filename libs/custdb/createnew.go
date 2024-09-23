package custdb

import (
	"fmt"
	"log"

	"github.com/transfashion/evoucher/libs/uniqid"
)

func (c *CustomerDB) CreateNew(phonenumber string, name string) (*Customer, error) {
	cust_id := uniqid.New(uniqid.Params{MoreEntropy: false})
	custaccess_id := cust_id

	// insert into mst_cust
	query := `
		INSERT INTO mst_cust (cust_id, cust_phone, cust_name, _createby) VALUES (?, ?, ?, '5effbb0a0f7d1')
	`
	_, err := c.Connection.Exec(query, cust_id, phonenumber, name)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// insert into mst_custaccess
	query = `
		INSERT INTO mst_custaccess (cust_id, custaccess_id, custaccesstype_id, custaccess_code, _createby) VALUES (?, ?, ?, ?, '5effbb0a0f7d1')
	`
	_, err = c.Connection.Exec(query, cust_id, custaccess_id, "WA", phonenumber)
	if err != nil {
		log.Println(err.Error())
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
