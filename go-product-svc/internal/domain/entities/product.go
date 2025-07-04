package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductBase struct {
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  uuid.UUID
	IsActive    bool
}

type Product struct {
	ID uuid.UUID
	ProductBase
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(product ProductBase) *Product {
	return &Product{
		ID:          uuid.New(),
		ProductBase: product,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (p *Product) UpdateDetails(product ProductBase) {
	p.Name = product.Name
	p.Description = product.Description
	p.Price = product.Price
	p.Stock = product.Stock
	p.CategoryID = product.CategoryID
	p.UpdatedAt = time.Now()
}

func (p *Product) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

func (p *Product) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

func (p *Product) UpdateStock(newStock int) {
	p.Stock = newStock
	p.UpdatedAt = time.Now()
}
