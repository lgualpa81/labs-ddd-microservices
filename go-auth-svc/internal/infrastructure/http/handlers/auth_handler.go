package handlers

import (
	"poc-auth-svc/internal/application/dtos"
	"poc-auth-svc/internal/application/usecases"
	"poc-auth-svc/internal/infrastructure/utils"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCase
	validator   *validator.Validate
}

func NewAuthHandler(authUseCase usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dtos.RegisterRequest
	if err := h.validateAndParseRequest(c, &req); err != nil {
		return err
	}

	response, err := h.authUseCase.Register(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered succesfully", response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dtos.LoginRequest
	if err := h.validateAndParseRequest(c, &req); err != nil {
		return err
	}

	response, err := h.authUseCase.Login(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", response)
}

func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	token, err := utils.ExtractBearerToken(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	response, err := h.authUseCase.ValidateToken(c.Context(), token)
	if err != nil || !response.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", fiber.Map{
			"valid": false,
		})
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Token is valid", response)
}

// validateAndParseRequest función genérica para validar Content-Type, parsear body y validar struct
func (h *AuthHandler) validateAndParseRequest(c *fiber.Ctx, req interface{}) error {
	// Validar Content-Type
	if err := utils.ValidateContentType(c, "application/json"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Parsear body
	if err := c.BodyParser(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", []string{err.Error()})
	}

	// Validar struct
	if err := h.validator.Struct(req); err != nil {
		validationErrors := utils.FormatValidationErrors(err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	return nil
}
