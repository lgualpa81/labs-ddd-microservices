package utils

import (
	"poc-product-svc/internal/application/dtos"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationService struct {
	validator *validator.Validate
}

func NewValidationService() *ValidationService {
	v := validator.New()

	// Registrar nombres de campos personalizados
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &ValidationService{validator: v}
}

func (v *ValidationService) ValidateStruct(s interface{}) []dtos.ValidationError {
	var validationErrors []dtos.ValidationError

	err := v.validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, dtos.ValidationError{
				Field:   err.Field(),
				Message: v.getErrorMessage(err),
				Value:   err.Value().(string),
			})
		}
	}

	return validationErrors
}

func (v *ValidationService) getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "This field must be at least " + err.Param() + " characters long"
	case "max":
		return "This field must be at most " + err.Param() + " characters long"
	case "gt":
		return "This field must be greater than " + err.Param()
	case "gte":
		return "This field must be greater than or equal to " + err.Param()
	case "lt":
		return "This field must be less than " + err.Param()
	case "lte":
		return "This field must be less than or equal to " + err.Param()
	case "email":
		return "This field must be a valid email address"
	case "uuid":
		return "This field must be a valid UUID"
	case "gtefield":
		return "This field must be greater than or equal to " + err.Param()
	default:
		return "This field is invalid"
	}
}

func (v *ValidationService) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}
