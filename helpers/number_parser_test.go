package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCountryCode(t *testing.T) {
	t.Run("ValidPhoneNumber", func(t *testing.T) {
		number := "+1 650-555-5555"
		expectedCode := "1"
		code, err := ParseCountryCode(number)
		assert.NoError(t, err)
		assert.Equal(t, expectedCode, code)
	})

	t.Run("ValidInternationalPhoneNumber", func(t *testing.T) {
		number := "+44 20 7123 1234"
		expectedCode := "44"
		code, err := ParseCountryCode(number)
		assert.NoError(t, err)
		assert.Equal(t, expectedCode, code)
	})

	t.Run("InvalidPhoneNumber", func(t *testing.T) {
		invalidNumber := "invalid-number"
		_, err := ParseCountryCode(invalidNumber)
		assert.Error(t, err)
	})
}
