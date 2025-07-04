package usecases

import (
	"context"
	"poc-product-svc/internal/application/dtos"
	"poc-product-svc/internal/domain/entities"
	"poc-product-svc/internal/domain/repositories"
	"poc-product-svc/internal/domain/services"

	"github.com/google/uuid"
)

type ProductUseCase struct {
	productRepo    repositories.ProductRepository
	productService *services.ProductService
}

func NewProductUseCase(
	productRepo repositories.ProductRepository,
	productService *services.ProductService) *ProductUseCase {
	return &ProductUseCase{
		productRepo:    productRepo,
		productService: productService,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, req dtos.CreateProductRequest) (*dtos.ProductResponse, error) {
	partialInput := entities.ProductBase{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	if err := uc.productService.ValidateProductCreation(ctx, partialInput); err != nil {
		return nil, err
	}
	input := entities.ProductBase{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}
	product := entities.NewProduct(input)
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return uc.entityToResponse(product), nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id uuid.UUID) (*dtos.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.entityToResponse(product), nil
}

func (uc *ProductUseCase) GetAllProducts(ctx context.Context, page, pageSize int) (*dtos.ProductListResponse, error) {
	offset := (page - 1) * pageSize
	products, err := uc.productRepo.GetAll(ctx, page, offset)
	if err != nil {
		return nil, err
	}
	var responses []dtos.ProductResponse
	for _, product := range products {
		responses = append(responses, *uc.entityToResponse(product))
	}
	totalPages := len(responses) / pageSize
	if len(responses)%pageSize > 0 {
		totalPages++
	}
	return &dtos.ProductListResponse{
		Products:   responses,
		Total:      len(responses),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req dtos.UpdateProductRequest) (*dtos.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	input := entities.ProductBase{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}
	product.UpdateDetails(input)
	if req.IsActive {
		product.Activate()
	} else {
		product.Deactivate()
	}
	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}
	return uc.entityToResponse(product), nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := uc.productRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc *ProductUseCase) UpdateStock(ctx context.Context, id uuid.UUID, req dtos.UpdateStockRequest) (*dtos.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	product.UpdateStock(req.Stock)
	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}
	return uc.entityToResponse(product), nil
}

func (uc *ProductUseCase) entityToResponse(product *entities.Product) *dtos.ProductResponse {
	return &dtos.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float64(product.Stock),
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
