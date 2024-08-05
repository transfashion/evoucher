package voucher

import "database/sql"

type VoucherDB struct {
	Connection *sql.DB
}

func NewVoucherDB(conn *sql.DB) *VoucherDB {
	return &VoucherDB{
		Connection: conn,
	}
}
