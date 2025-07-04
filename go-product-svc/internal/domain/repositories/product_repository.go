package repositories

import (
	"context"
	"poc-product-svc/internal/domain/entities"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error)
	SearchByName(ctx context.Context, name string) ([]*entities.Product, error)
}

type ProductCacheRepository interface {
	Set(ctx context.Context, key string, product *entities.Product) error
	Get(ctx context.Context, key string) (*entities.Product, error)
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}
