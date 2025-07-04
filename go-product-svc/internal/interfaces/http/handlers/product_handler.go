package handlers

import (
	"poc-product-svc/internal/application/dtos"
	"poc-product-svc/internal/application/usecases"
	"poc-product-svc/internal/shared/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productUC *usecases.ProductUseCase
	validator *utils.ValidationService
}

func NewProductHandler(
	productUC *usecases.ProductUseCase,
	validator *utils.ValidationService,
) *ProductHandler {
	return &ProductHandler{
		productUC: productUC,
		validator: validator,
	}
}

func (c *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	var req dtos.CreateProductRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid request body", err.Error()),
		)
	}

	if validationErrors := c.validator.ValidateStruct(req); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewValidationErrorResponse(validationErrors),
		)
	}

	product, err := c.productUC.CreateProduct(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse("failed to create product", err.Error()),
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(
		dtos.NewSuccessResponse("product created successfully", product),
	)

}

func (c *ProductHandler) GetProduct(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid product ID", "Product ID must be a valid UUID"),
		)
	}
	product, err := c.productUC.GetProductByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(
			dtos.NewErrorResponse("Product not found", err.Error()),
		)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		dtos.NewSuccessResponse("Product retrieved successfully", product),
	)
}

func (c *ProductHandler) GetAllProducts(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, err := c.productUC.GetAllProducts(ctx.Context(), page, pageSize)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse("failed to retrieve products", err.Error()),
		)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		dtos.NewSuccessResponse("products retrieve successfully", products),
	)
}

func (c *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid product ID", "Product ID must be a valid UUID"),
		)
	}
	var req dtos.UpdateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid request body", err.Error()),
		)
	}

	if validationErrors := c.validator.ValidateStruct(req); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewValidationErrorResponse(validationErrors),
		)
	}
	product, err := c.productUC.UpdateProduct(ctx.Context(), id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse("Failed to update product", err.Error()),
		)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		dtos.NewSuccessResponse("Product updated successfully", product),
	)
}

func (c *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid product ID", "Product ID must be a valid UUID"),
		)
	}
	err = c.productUC.DeleteProduct(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse("Failed to delete product", err.Error()),
		)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		dtos.NewSuccessResponse("Product deleted successfully", nil),
	)
}

func (c *ProductHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid product ID", "Product ID must be a valid UUID"),
		)
	}

	var req dtos.UpdateStockRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse("Invalid request body", err.Error()),
		)
	}

	if validationErrors := c.validator.ValidateStruct(req); len(validationErrors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dtos.NewValidationErrorResponse(validationErrors),
		)
	}

	product, err := c.productUC.UpdateStock(ctx.Context(), id, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse("Failed to update stock", err.Error()),
		)
	}
	return ctx.Status(fiber.StatusOK).JSON(
		dtos.NewSuccessResponse("Stock updated successfully", product),
	)
}
