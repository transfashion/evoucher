package voucher

import (
	"database/sql"
	"log"

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
		log.Println("data voucher tidak ditemukan")
		return nil, nil
	} else if err != nil {
		log.Println(err.Error())
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

func (v *VoucherDB) GetAllUnusedVoucherByPhoneNumber(voubatch_id string, phonenumber string) (vouchers []*Voucher, err error) {

	query := `
		select 
		vou_id, vou_isactive
		from
		mst_vou
		where
		voubatch_id=? and vou_assignto=? and vou_isuse=0
	`

	rows, err := v.Connection.Query(query, voubatch_id, phonenumber)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	vouchers = make([]*Voucher, 0)
	for rows.Next() {
		var vou_id *string
		var vou_isactive bool

		vou_id = new(string)
		err := rows.Scan(&vou_id, &vou_isactive)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		vouchers = append(vouchers, &Voucher{
			Id:       helper.IsStringNil(vou_id, ""),
			IsActive: vou_isactive,
			IsUse:    false,
			UsedBy:   "",
		})
	}

	return vouchers, nil
}

func (v *VoucherDB) GetSingleVoucherByPhoneNumber(voubatch_id string, phonenumber string) (voucher *Voucher, err error) {
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
		voubatch_id=? and vou_assignto=?
	`

	row := v.Connection.QueryRow(query, voubatch_id, phonenumber)
	err = row.Scan(&vou_id, &vou_isactive, &vou_isuse, &vou_useby, &vou_usedate)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err.Error())
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
