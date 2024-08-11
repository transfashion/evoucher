package voucher

import (
	"database/sql"
	"fmt"
	"time"
)

type VoucherBatch struct {
	Type        *string
	BrandId     *string
	Description *string
	Greeting    *string
	CrmEventId  *string
	DtStart     time.Time
	DtEnd       time.Time
	DtActive    time.Time
	DtExpired   time.Time
	Value       float32

	_dtStart   *string
	_dtEnd     *string
	_dtActive  *string
	_dtExpired *string
}

func (v *VoucherDB) GetVoucherBatch(voubatch_id string) (*VoucherBatch, error) {

	query := `
		select 
		voutype_id,brand_id,voubatch_descr,voubatch_greeting,crmevent_id,voubatch_dtstart,voubatch_dtend,voubatch_dtactive,voubatch_dtexpired,voubatch_value
		from mst_voubatch where voubatch_id = ?
	`

	vbatch := &VoucherBatch{}
	row := v.Connection.QueryRow(query, voubatch_id)
	err := row.Scan(
		&vbatch.Type,
		&vbatch.BrandId,
		&vbatch.Description,
		&vbatch.Greeting,
		&vbatch.CrmEventId,
		&vbatch._dtStart,
		&vbatch._dtEnd,
		&vbatch._dtActive,
		&vbatch._dtExpired,
		&vbatch.Value,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("voucher batch not found")
	} else if err != nil {
		return nil, err
	}

	vbatch.DtStart = toTime(*vbatch._dtStart)
	vbatch.DtEnd = toTime(*vbatch._dtEnd)
	vbatch.DtActive = toTime(*vbatch._dtActive)
	vbatch.DtExpired = toTime(*vbatch._dtExpired)

	return vbatch, nil
}

func toTime(tanggalStr string) time.Time {
	layout := "2006-01-02"
	tanggal, _ := time.Parse(layout, tanggalStr)
	return tanggal
}
