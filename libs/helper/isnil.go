package helper

func IsStringNil(str *string, valueifnil string) string {
	if str == nil {
		return valueifnil
	}
	return *str
}
