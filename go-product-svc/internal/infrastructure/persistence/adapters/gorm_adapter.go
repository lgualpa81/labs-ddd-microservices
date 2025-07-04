package adapters

import (
	"context"
	"fmt"
	"poc-product-svc/internal/domain/entities"
	"poc-product-svc/internal/infrastructure/persistence/postgres/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormAdapter implementa ORMAdapter usando GORM
type GormAdapter struct {
	db *gorm.DB
}

// NewGormAdapter crea una nueva instancia del adapter GORM
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

func (g *GormAdapter) Create(ctx context.Context, product *entities.Product) error {
	model := models.NewProductModelFromEntity(product)

	if err := g.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// Actualizar la entidad con los datos generados por la base de datos
	*product = *model.ToEntity()
	return nil
}

func (g *GormAdapter) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	var model models.ProductModel

	err := g.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return model.ToEntity(), nil
}

func (g *GormAdapter) GetAll(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	var models []models.ProductModel

	query := g.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	products := make([]*entities.Product, len(models))
	for i, model := range models {
		products[i] = model.ToEntity()
	}

	return products, nil
}

func (g *GormAdapter) Update(ctx context.Context, product *entities.Product) error {
	model := models.NewProductModelFromEntity(product)
	model.UpdatedAt = time.Now().UTC()

	result := g.db.WithContext(ctx).
		Where("id = ?", product.ID).
		Updates(model)

	if result.Error != nil {
		return fmt.Errorf("failed to update product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", product.ID)
	}

	// Actualizar la entidad con el nuevo timestamp
	product.UpdatedAt = model.UpdatedAt
	return nil
}

func (g *GormAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	result := g.db.WithContext(ctx).Where("id = ?", id).Delete(&models.ProductModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product with id %s not found", id)
	}

	return nil
}

func (g *GormAdapter) GetByCategory(ctx context.Context, categoryID uuid.UUID) ([]*entities.Product, error) {
	var models []models.ProductModel

	err := g.db.WithContext(ctx).
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	products := make([]*entities.Product, len(models))
	for i, model := range models {
		products[i] = model.ToEntity()
	}

	return products, nil
}

func (g *GormAdapter) SearchByName(ctx context.Context, name string) ([]*entities.Product, error) {
	var models []models.ProductModel

	searchPattern := "%" + name + "%"
	err := g.db.WithContext(ctx).
		Where("name ILIKE ? AND is_active = ?", searchPattern, true).
		Order("name ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search products by name: %w", err)
	}

	products := make([]*entities.Product, len(models))
	for i, model := range models {
		products[i] = model.ToEntity()
	}

	return products, nil
}

func (g *GormAdapter) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// GormTransactionAdapter implementa TransactionAdapter
type GormTransactionAdapter struct {
	tx      *gorm.DB
	adapter *GormAdapter
}

func (g *GormTransactionAdapter) Commit() error {
	if err := g.tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (g *GormTransactionAdapter) Rollback() error {
	if err := g.tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (g *GormTransactionAdapter) GetORMAdapter() ORMAdapter {
	return g.adapter
}

// GormDatabaseAdapter implementa DatabaseAdapter
type GormDatabaseAdapter struct {
	db *gorm.DB
}

func NewGormDatabaseAdapter(db *gorm.DB) *GormDatabaseAdapter {
	return &GormDatabaseAdapter{db: db}
}

func (g *GormDatabaseAdapter) BeginTransaction(ctx context.Context) (TransactionAdapter, error) {
	tx := g.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	return &GormTransactionAdapter{
		tx:      tx,
		adapter: NewGormAdapter(tx),
	}, nil
}

func (g *GormDatabaseAdapter) Ping(ctx context.Context) error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

func (g *GormDatabaseAdapter) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
