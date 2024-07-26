package models

import (
	"github.com/transfashion/evoucher/custdb"
	"github.com/transfashion/evoucher/qiscus"
)

type ApplicationData struct {
	Qiscus     *qiscus.Qiscus
	CustomerDb *custdb.CustomerDB
}
