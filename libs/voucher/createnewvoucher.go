package voucher

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/transfashion/evoucher/libs/uniqid"
)

func (v *VoucherDB) CreateNewVoucher(voubatch_id string, phonenumber string, customername string, code string) (voucher *Voucher, err error) {
	temp_vou_id := uniqid.New(uniqid.Params{MoreEntropy: true})

	var voutype_id, voubatch_code *string
	var vou_no, rndmin, rndmax int
	var vou_value float32

	voutype_id = new(string)
	voubatch_code = new(string)
	// vou_no,vou_value,voutype_id,rndmin,rndmax,voubatch_code

	// buat draft voucher
	log.Println("create tempvoucher", temp_vou_id, voubatch_id)
	query := "call vou_create(?, ?)"
	row := v.Connection.QueryRow(query, temp_vou_id, voubatch_id)
	err = row.Scan(&vou_no, &vou_value, &voutype_id, &rndmin, &rndmax, &voubatch_code)
	if err == sql.ErrNoRows {
		log.Println("error saat membuat draft voucher")
		return nil, fmt.Errorf("error saat membuat draft voucher")
	} else if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if vou_no > 9999 {
		log.Println("jumlah quota voucher sudah mencapai batas")
		return nil, fmt.Errorf("jumlah quota voucher sudah mencapai batas")
	}

	// compose protitipe voucher
	voucher = &Voucher{
		No:        vou_no,
		Value:     vou_value,
		Rmin:      rndmin,
		Rmax:      rndmax,
		BatchCode: *voubatch_code,
		Type:      *voutype_id,
	}

	// buat random number dari rndmin sampai rndmax
	ran := rand.Intn(rndmax-rndmin+1) + rndmin

	// compose voucher code
	i_voubatch_id, _ := strconv.Atoi(voubatch_id)
	i_voubatch_code, _ := strconv.Atoi(*voubatch_code)
	t := float64(vou_no + i_voubatch_id + i_voubatch_code)
	b := float64(ran)
	n := math.Floor(t / b)
	p := int(n) % ran
	parstr := fmt.Sprintf("%02d", p)
	parity := parstr[len(parstr)-2:]
	nopad := fmt.Sprintf("%04d", vou_no)
	voucher_id := fmt.Sprintf("%s%s%s%d%s", voubatch_id, *voubatch_code, nopad, ran, parity)

	// update draft voucher ke voucher yang telah terbentuk
	query = `
		update mst_vou 
		set
		vou_id = ?,
		vou_ran = ?,
		vou_parity = ?,
		vou_infocode = 'WA',
		vou_assigncode = ?,
		vou_assignto = ?,
		vou_assigntoname = ?,
		vou_isactive = 1
		where
		vou_id = ?
	`

	log.Println("updating tempvoucher", temp_vou_id, voucher_id)
	_, err = v.Connection.Exec(query, voucher_id, ran, parity, code, phonenumber, customername, temp_vou_id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	voucher.Id = voucher_id

	vou, err := v.GetVoucher(voucher.Id)
	if err != nil {
		log.Println("gagal mendapatkan kembali data voucher")
		return nil, err
	}

	return vou, nil
}
