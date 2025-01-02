package utils

import (
	"fmt"
	database "goon/db"
	"log"

	"github.com/go-playground/validator/v10"
)

func IsUniqueValue(tableName, fieldName, value string) bool {
	var count int64

	result := database.DB.Table(tableName).Where(fieldName+" = ?", value).Count(&count)

	if result.Error != nil {
		log.Println("Error: ", result.Error)
		return false
	}

	return count > 0
}

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)

	for _, err := range errs {
		fmt.Println()
		switch err.Tag() {
		case "required":
			errorMessages[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be a valid email address", err.Field())
		}
	}
}
