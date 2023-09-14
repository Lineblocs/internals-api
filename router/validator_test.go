package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {

	t.Run("Should validate without errors", func(t *testing.T) {

		validator := NewValidator()

		type TestStruct struct {
			Field1 string `validate:"required"`
			Field2 int    `validate:"gte=10"`
		}

		testData := TestStruct{
			Field1: "ValidValue",
			Field2: 15,
		}

		err := validator.Validate(testData)
		assert.NoError(t, err)
	})

	t.Run("Should return error for incorrect struct", func(t *testing.T) {

		validator := NewValidator()

		type TestStruct struct {
			Field1 string `validate:"required"`
			Field2 int    `validate:"gte=10"`
		}

		invalidData := TestStruct{
			Field1: "",
			Field2: 5,
		}

		err := validator.Validate(invalidData)
		assert.Error(t, err)
	})
}
