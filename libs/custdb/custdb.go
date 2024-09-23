package custdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerDB struct {
	Connection *sql.DB
}

type Customer struct {
	Id          string
	PhoneNumber string
	Name        string
	Gender      string
}

func NewCustomerDB(conn *sql.DB) *CustomerDB {
	return &CustomerDB{
		Connection: conn,
	}
}
