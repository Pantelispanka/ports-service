package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

// The port struct, and a validation of each field.
type Port struct {
	Unloc       string    `json:"unloc" validate:"required"`
	Name        string    `json:"name"`
	Coordinates []float64 `json:"coordinates" validate:"dive,latitude,longitude"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	Allias      []string  `json:"allias"`
	Regions     []string  `json:"regions"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code" validate:"required,len=5"`
}

// Validate the Port struct.
// If the validation fails and error will be returned with the details of the error.
func (p *Port) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		var validationErrors []string
		for _, vErr := range err.(validator.ValidationErrors) {
			errMessage := fmt.Sprintf("'%s' has a value of '%v' which does not satisfy '%s'.\n", vErr.Field(), vErr.Value(), vErr.Tag())
			validationErrors = append(validationErrors, errMessage)
		}
		return fmt.Errorf("error in struct: %+v. details: '%v'", *p, strings.Join(validationErrors, ", "))
	}
	return nil
}
