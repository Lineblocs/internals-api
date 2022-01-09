package helpers



import (
	"github.com/ttacon/libphonenumber"
	"strconv"
)


func ParseCountryCode( number string ) (string, error) {
	num, err := libphonenumber.Parse(number, "")

	if err != nil {
		return "", err
	}

	code := num.GetCountryCode()

	return strconv.Itoa( int( code ) ), nil
}