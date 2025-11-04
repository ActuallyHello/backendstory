package common

import "github.com/go-playground/validator/v10"

func CollectValidationDetails(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)

	details := make(map[string]string)
	if len(validationErrors) > 0 {
		for _, validationError := range validationErrors {
			details[validationError.Field()] = validationError.Tag()
		}
	}

	return details
}
