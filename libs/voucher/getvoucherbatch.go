package voucher

import (
	"database/sql"
	"fmt"
	"time"
)

type VoucherBatch struct {
	Id string

	Type              *string
	BrandId           *string
	Description       *string
	Greeting          *string
	CrmEventId        *string
	DtStart           time.Time
	DtEnd             time.Time
	DtActive          time.Time
	DtExpired         time.Time
	Value             float32
	IsGenerated       bool
	IsMultiple        bool
	UseActivationCode bool

	_dtStart   *string
	_dtEnd     *string
	_dtActive  *string
	_dtExpired *string
}

func (v *VoucherDB) GetVoucherBatch(voubatch_id string) (*VoucherBatch, error) {

	query := `
		select 
			voutype_id,brand_id,voubatch_descr,voubatch_greeting,crmevent_id,
			voubatch_dtstart,voubatch_dtend,
			voubatch_dtactive,voubatch_dtexpired,
			voubatch_value,
			voubatch_isusecodeact,
			voubatch_isgenerate
		from 
			mst_voubatch 
		where 
			voubatch_id = ?
	`

	var isusecodeact bool

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
		&isusecodeact,
		&vbatch.IsGenerated,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("voucher batch not found")
	} else if err != nil {
		return nil, err
	}

	vbatch.Id = voubatch_id
	vbatch.DtStart = toTime(*vbatch._dtStart)
	vbatch.DtEnd = toTime(*vbatch._dtEnd)
	vbatch.DtActive = toTime(*vbatch._dtActive)
	vbatch.DtExpired = toTime(*vbatch._dtExpired)

	vbatch.UseActivationCode = isusecodeact
	if isusecodeact {
		vbatch.IsMultiple = true
	} else {
		vbatch.IsMultiple = false
	}

	return vbatch, nil
}

func toTime(tanggalStr string) time.Time {
	layout := "2006-01-02"
	tanggal, _ := time.Parse(layout, tanggalStr)
	return tanggal
}
