package models

import (
	"poc-product-svc/internal/domain/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductModel representa el modelo de datos para GORM
type ProductModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(255);not null;index"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(10,2);not null;check:price >= 0"`
	Stock       int       `gorm:"not null;check:stock >= 0"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;index"`
	IsActive    bool      `gorm:"default:true;index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// TableName especifica el nombre de la tabla
func (ProductModel) TableName() string {
	return "products"
}

// BeforeCreate hook de GORM para generar UUID antes de crear
func (p *ProductModel) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// ToEntity convierte el modelo GORM a entidad de dominio
func (p *ProductModel) ToEntity() *entities.Product {
	return &entities.Product{
		ID: p.ID,
		ProductBase: entities.ProductBase{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			CategoryID:  p.CategoryID,
			IsActive:    p.IsActive,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo GORM
func (p *ProductModel) FromEntity(product *entities.Product) {
	p.ID = product.ID
	p.Name = product.Name
	p.Description = product.Description
	p.Price = product.Price
	p.Stock = product.Stock
	p.CategoryID = product.CategoryID
	p.IsActive = product.IsActive
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
}

// NewProductModelFromEntity crea un nuevo modelo desde una entidad
func NewProductModelFromEntity(product *entities.Product) *ProductModel {
	model := &ProductModel{}
	model.FromEntity(product)
	return model
}
