package services

import (
	"context"
	"poc-product-svc/internal/domain/entities"
	"poc-product-svc/internal/domain/errors"
	"poc-product-svc/internal/domain/repositories"
)

type ProductService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) ValidateProductCreation(ctx context.Context, product entities.ProductBase) error {
	if product.Name == "" {
		return errors.ErrInvalidProductName
	}
	if product.Price < 0 {
		return errors.ErrInvalidPrice
	}
	if product.Stock < 0 {
		return errors.ErrInvalidStock
	}
	products, err := s.productRepo.SearchByName(ctx, product.Name)
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return errors.ErrProductAlreadyExists
	}
	return nil
}

func ValidateStockUpdate(currentStock, requestQuantity int) error {
	if currentStock < requestQuantity {
		return errors.ErrInsufficientStock
	}
	return nil
}
