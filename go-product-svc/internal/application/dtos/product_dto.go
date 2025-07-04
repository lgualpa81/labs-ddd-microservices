package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Description string    `json:"description" validate:"required,min=5,max=500"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Description string    `json:"description" validate:"required,min=5,max=500"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	IsActive    bool      `json:"is_active"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uuid.UUID `json:"category_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,gte=0"`
}

type ProductSearchRequest struct {
	Name       string    `json:"name,omitempty"`
	CategoryID uuid.UUID `json:"category_id,omitempty"`
	MinPrice   float64   `json:"min_price,omitempty" validate:"gte=0"`
	MaxPrice   float64   `json:"max_price,omitempty" validate:"gtefield=MinPrice"`
	Page       int       `json:"page" validate:"min=1"`
	PageSize   int       `json:"page_size" validate:"min=1,max=100"`
}
