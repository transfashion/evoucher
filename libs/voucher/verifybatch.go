package voucher

import (
	"fmt"
	"time"
)

func (v *VoucherDB) VerifyBatch(vb *VoucherBatch) (bool, error) {
	currentTime := time.Now()
	format := "2006-01-02"
	dtf := currentTime.Format(format)
	now, _ := time.Parse(format, dtf)

	dtst := vb.DtStart.Format(format)
	dten := vb.DtEnd.Format(format)
	if vb.DtStart.After(now) {
		return false, fmt.Errorf("periode voucher belum dimulai, voucher bisa direquest mulai tanggal %s s/d %s", dtst, dten)
	}

	if vb.DtEnd.Before(now) {
		return false, fmt.Errorf("periode voucher sudah habis per tanggal %s", dten)
	}

	return true, nil
}
