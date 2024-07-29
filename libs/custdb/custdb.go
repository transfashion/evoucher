package custdb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerDB struct {
	// ...
}

type DatabaseConfig struct {
	Server   string `yaml:"server"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Customer struct {
	PhoneNumber string
	Name        string
}

func NewCustomerDB(dbconf *DatabaseConfig) *CustomerDB {
	ConnectDatabase(dbconf.Server, dbconf.Name, dbconf.Username, dbconf.Password)
	return &CustomerDB{}
}

func ConnectDatabase(host string, dbname string, username string, password string) (*sql.DB, error) {
	server := fmt.Sprintf("%s:%d", host, 3306)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, server, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
