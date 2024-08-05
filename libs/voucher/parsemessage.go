package voucher

import (
	"regexp"
)

type VoucherMessageIntent struct {
	Ref string
}

func (v *VoucherDB) ParseMessage(msg string) *VoucherMessageIntent {

	regex := regexp.MustCompile(`\[ref:(.*?)\]`)
	matches := regex.FindStringSubmatch(msg)
	if len(matches) > 1 {
		ref := matches[1]
		return &VoucherMessageIntent{
			Ref: ref,
		}
	} else {
		return nil
	}
}
