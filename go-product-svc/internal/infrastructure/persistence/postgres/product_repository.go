package postgres

import (
	"context"
	"poc-product-svc/internal/domain/entities"
	"poc-product-svc/internal/domain/repositories"
	"poc-product-svc/internal/infrastructure/persistence/adapters"

	"github.com/google/uuid"
)

// // ProductRepository implementa el repositorio de productos usando el adapter ORM
type ProductRepository struct {
	adapter adapters.ORMAdapter
}

// NewProductRepository crea una nueva instancia del repositorio
func NewProductRepository(adapter adapters.ORMAdapter) repositories.ProductRepository {
	return &ProductRepository{
		adapter: adapter,
	}
}

// Create implements repositories.ProductRepository.
func (p *ProductRepository) Create(ctx context.Context, product *entities.Product) error {
	return p.adapter.Create(ctx, product)
}

// Delete implements repositories.ProductRepository.
func (p *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return p.adapter.Delete(ctx, id)
}

// GetAll implements repositories.ProductRepository.
func (p *ProductRepository) GetAll(ctx context.Context, limit int, offset int) ([]*entities.Product, error) {
	return p.adapter.GetAll(ctx, limit, offset)
}

// GetByCategory implements repositories.ProductRepository.
func (p *ProductRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error) {
	return p.adapter.GetByCategory(ctx, categoryID)
}

// GetByID implements repositories.ProductRepository.
func (p *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	return p.adapter.GetByID(ctx, id)
}

// SearchByName implements repositories.ProductRepository.
func (p *ProductRepository) SearchByName(ctx context.Context, name string) ([]*entities.Product, error) {
	return p.adapter.SearchByName(ctx, name)
}

// Update implements repositories.ProductRepository.
func (p *ProductRepository) Update(ctx context.Context, product *entities.Product) error {
	return p.adapter.Update(ctx, product)
}
