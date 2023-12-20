package domain_test

import (
	"ports-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortValidation(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		port := domain.Port{
			Unloc: "AEAUH",
			City:  "Abu Dhabi",
			Code:  "52001",
		}
		err := port.Validate()
		assert.NoError(t, err)
	})

	t.Run("pass with coordinates", func(t *testing.T) {
		coordinates := []float64{12.876, -42.12}
		port := domain.Port{
			Unloc:       "AEAUH",
			City:        "Abu Dhabi",
			Code:        "52001",
			Coordinates: coordinates,
		}
		err := port.Validate()
		assert.NoError(t, err)
	})

	t.Run("pass wrong latitude", func(t *testing.T) {
		errMessage := "''Coordinates[1]' has a value of '120000' which does not satisfy 'latitude'"
		coordinates := []float64{0.2, 120000}
		port := domain.Port{
			Unloc:       "AEAUH",
			City:        "Abu Dhabi",
			Code:        "52001",
			Coordinates: coordinates,
		}
		err := port.Validate()
		assert.ErrorContains(t, err, errMessage)
	})

	t.Run("A correct code", func(t *testing.T) {
		port := domain.Port{
			Unloc: "AEAUH",
			City:  "Abu Dhabi",
			Code:  "30000",
		}
		err := port.Validate()
		assert.NoError(t, err)
	})

	t.Run("A wrong code", func(t *testing.T) {
		errorMessage := "''Code' has a value of '300000' which does not satisfy 'len'"
		port := domain.Port{
			Unloc: "AEAUH",
			City:  "Abu Dhabi",
			Code:  "300000",
		}
		err := port.Validate()
		assert.ErrorContains(t, err, errorMessage)
	})

}
