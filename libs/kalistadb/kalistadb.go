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
	Port     int    `yaml:"port"`
}

var conn *sql.DB

func ConnectDatabase(host string, dbname string, port int, username string, password string) error {
	server := fmt.Sprintf("%s:%d", host, port)
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
