package services

import (
	"backend/internal/model"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateRequestData(data *model.RequestData) error {
	// Create a new validator instance
	validator := validator.New()

	// Define the validation rules
	err := validator.Struct(data)
	if err != nil {
		// Return a descriptive error message if validation fails
		return fmt.Errorf("invalid request data: %s", err.Error())
	}

	return nil
}
