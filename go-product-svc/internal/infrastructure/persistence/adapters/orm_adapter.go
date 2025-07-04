package adapters

import (
	"context"
	"poc-product-svc/internal/domain/entities"

	"github.com/google/uuid"
)

// ORMAdapter define la interfaz para adapters de ORM
type ORMAdapter interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error)
	SearchByName(ctx context.Context, name string) ([]*entities.Product, error)
	Close() error
}

// DatabaseAdapter define operaciones genéricas de base de datos
type DatabaseAdapter interface {
	BeginTransaction(ctx context.Context) (TransactionAdapter, error)
	Ping(ctx context.Context) error
	Close() error
}

// TransactionAdapter define operaciones de transacción
type TransactionAdapter interface {
	Commit() error
	Rollback() error
	GetORMAdapter() ORMAdapter
}
