package voucher

import (
	"regexp"
)

type VoucherMessageIntent struct {
	VoubatchId string
}

func (v *VoucherDB) ParseMessage(msg string) *VoucherMessageIntent {

	regex := regexp.MustCompile(`\[ref:(.*?)\]`)
	matches := regex.FindStringSubmatch(msg)
	if len(matches) > 1 {
		ref := matches[1]
		return &VoucherMessageIntent{
			VoubatchId: ref,
		}
	} else {
		return nil
	}
}
