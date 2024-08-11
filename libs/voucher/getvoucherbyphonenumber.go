package voucher

import (
	"database/sql"

	"github.com/transfashion/evoucher/libs/helper"
)

func (v *VoucherDB) GetVoucherByPhoneNumber(voubatch_id string, phonenumber string, code string) (voucher *Voucher, err error) {
	var vou_id, vou_useby, vou_usedate *string
	var vou_isactive, vou_isuse bool

	vou_id = new(string)
	vou_useby = new(string)
	vou_usedate = new(string)

	query := `
		select 
		vou_id, vou_isactive, vou_isuse, vou_useby, vou_usedate
		from
		mst_vou
		where
		voubatch_id=? and vou_assignto=? and vou_assigncode=?
	`

	row := v.Connection.QueryRow(query, voubatch_id, phonenumber, code)
	err = row.Scan(&vou_id, &vou_isactive, &vou_isuse, &vou_useby, &vou_usedate)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	voucher = &Voucher{
		Id:       helper.IsStringNil(vou_id, ""),
		IsActive: vou_isactive,
		IsUse:    vou_isuse,
		UsedBy:   helper.IsStringNil(vou_useby, ""),
	}

	return voucher, nil
}
