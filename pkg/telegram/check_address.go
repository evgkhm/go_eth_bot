package telegram

import "regexp"

// IsValidAddress функция проверки валидности eth адреса
func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
	return re.MatchString(v)
}

//func IsValidAddressFromMap(v int64) bool {
//	re := regexp.MustCompile("^0x[\\da-fA-F]{40}$")
//	return re.MatchString(v)
//}
