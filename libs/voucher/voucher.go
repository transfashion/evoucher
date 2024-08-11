package voucher

import (
	"database/sql"
	"time"
)

type VoucherDB struct {
	Connection *sql.DB
}

type Voucher struct {
	Id             string
	BatchId        string
	IsActive       bool
	IsUse          bool
	AssignTo       string
	AssignName     string
	AssignCode     string
	UsedBy         string
	UsedDate       string
	ExpiredDate    time.Time
	Description    string
	BrandId        string
	HeaderLogoData []byte
	IconLogoData   []byte

	No        int
	Rmin      int
	Rmax      int
	Value     float32
	BatchCode string
	Type      string
}

func NewVoucherDB(conn *sql.DB) *VoucherDB {
	return &VoucherDB{
		Connection: conn,
	}
}
