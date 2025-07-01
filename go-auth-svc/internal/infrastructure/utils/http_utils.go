package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// StandardResponse estructura estandarizada para todas las respuestas
type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SuccessResponse respuesta estandarizada para casos exitosos
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := StandardResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	return c.Status(status).JSON(response)
}

// ErrorResponse respuesta estandarizada para errores
func ErrorResponse(c *fiber.Ctx, status int, message string, details interface{}) error {
	response := StandardResponse{
		Success:   false,
		Message:   "Request failed",
		Error:     message,
		Details:   details,
		Timestamp: time.Now(),
	}

	return c.Status(status).JSON(response)
}

// ValidateContentType valida que el Content-Type sea el esperado
func ValidateContentType(c *fiber.Ctx, expectedType string) error {
	contentType := c.Get("Content-Type")
	if !strings.Contains(contentType, expectedType) {
		return fmt.Errorf("Content-Type must be %s", expectedType)
	}
	return nil
}

// ExtractBearerToken extrae el token del header Authorization
func ExtractBearerToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", errors.New("bearer token required")
	}

	return tokenString, nil
}

// FormatValidationErrors convierte errores de validación en mensajes legibles
func FormatValidationErrors(err error) []string {
	validationErrors := make([]string, 0)

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErrors {
			validationErrors = append(validationErrors, formatValidationError(fieldError))
		}
	}

	return validationErrors
}

// formatValidationError formatea un error de validación individual
func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
