package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

var ErrValidation = fmt.Errorf("validation error")

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

func ValidateDepartmentName(name string) error {
	if err := validate.Var(name, "required,min=1,max=200"); err != nil {
		return fmt.Errorf("%w: department name must be between 1 and 200 characters", ErrValidation)
	}
	return nil
}

func ValidateEmployeeName(name string) error {
	if err := validate.Var(name, "required,min=1,max=200"); err != nil {
		return fmt.Errorf("%w: employee name must be between 1 and 200 characters", ErrValidation)
	}
	return nil
}

func ValidateEmployeePosition(pos string) error {
	if err := validate.Var(pos, "required,min=1,max=200"); err != nil {
		return fmt.Errorf("%w: employee position must be between 1 and 200 characters", ErrValidation)
	}
	return nil
}
