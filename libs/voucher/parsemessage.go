package voucher

import (
	"fmt"
	"regexp"
)

type VoucherMessageIntent struct{}

func (v *VoucherDB) ParseMessage(msg string) *VoucherMessageIntent {
	regex := regexp.MustCompile(`#\w+`)
	matches := regex.FindAllString(msg, -1)

	// regex := regexp.MustCompile(`\[(.*?)\]`)
	// matches := regex.FindAllStringSubmatch(text, -1)

	fmt.Println(matches)

	return nil
}
