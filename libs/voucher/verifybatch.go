package voucher

import (
	"fmt"
	"time"
)

func (v *VoucherDB) VerifyBatchRequest(vb *VoucherBatch) (bool, error) {

	// cek apakah voucher sudah digenerate
	if !vb.IsGenerated {
		return false, fmt.Errorf("maaf, kode batch voucher tidak valid. Request voucher tidak dapat diproses")
	}

	// Cek periode voucher bisa di generate
	currentTime := time.Now()
	format := "2006-01-02"
	dtf := currentTime.Format(format)
	now, _ := time.Parse(format, dtf)
	dtst := vb.DtStart.Format(format)
	dten := vb.DtEnd.Format(format)
	if vb.DtStart.After(now) {
		return false, fmt.Errorf("maaf, periode voucher belum dimulai, Saat ini request voucher belum dapat diproses. voucher bisa direquest mulai tanggal %s s/d %s", dtst, dten)
	}

	if vb.DtEnd.Before(now) {
		return false, fmt.Errorf("maaf, periode voucher sudah habis per tanggal %s. Request voucher tidak dapat diproses", dten)
	}

	return true, nil
}
