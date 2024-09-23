package voucher

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/transfashion/evoucher/libs/helper"
)

func (v *VoucherDB) GetVoucher(vou_id string) (*Voucher, error) {

	var brand_id, vou_dtexpired, voubatch_descr, voubatch_id *string
	var vou_value float32
	var vou_isactive, vou_isuse bool

	brand_id = new(string)
	vou_dtexpired = new(string)
	voubatch_descr = new(string)
	voubatch_id = new(string)

	query := `
		select 
		A.vou_value , A.vou_isactive , A.vou_isuse , A.vou_dtexpired,
		B.brand_id, B.voubatch_descr, B.voubatch_id 
		from mst_vou A inner join mst_voubatch B on B.voubatch_id =A.voubatch_id 
		WHERE 
		vou_id = ?
	`

	row := v.Connection.QueryRow(query, vou_id)
	err := row.Scan(
		&vou_value,
		&vou_isactive,
		&vou_isuse,
		&vou_dtexpired,
		&brand_id,
		&voubatch_descr,
		&voubatch_id,
	)
	if err == sql.ErrNoRows {
		log.Println("voucher batch not found")
		return nil, fmt.Errorf("voucher batch not found")
	} else if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	voucher := &Voucher{}
	voucher.Id = vou_id
	voucher.IsActive = vou_isactive
	voucher.IsUse = vou_isuse
	voucher.ExpiredDate = toTime(*vou_dtexpired)
	voucher.Value = vou_value
	voucher.Description = helper.IsStringNil(voubatch_descr, "")
	voucher.BrandId = helper.IsStringNil(brand_id, "")
	voucher.BatchId = helper.IsStringNil(voubatch_id, "")

	return voucher, nil
}
