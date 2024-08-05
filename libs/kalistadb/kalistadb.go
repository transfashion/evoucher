package kalistadb

import (
	"database/sql"
	"fmt"
)

type DatabaseConfig struct {
	Server   string `yaml:"server"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var conn *sql.DB

func ConnectDatabase(host string, dbname string, username string, password string) error {
	server := fmt.Sprintf("%s:%d", host, 3306)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, server, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	conn = db
	return nil
}

func GetConnection() *sql.DB {
	return conn
}
