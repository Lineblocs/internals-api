package helpers

import (
	"strconv"

	"github.com/ttacon/libphonenumber"
)

func ParseCountryCode(number string) (string, error) {
	num, err := libphonenumber.Parse(number, "")

	if err != nil {
		return "", err
	}

	code := num.GetCountryCode()

	return strconv.Itoa(int(code)), nil
}
