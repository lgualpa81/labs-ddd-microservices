package errors

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInvalidProductName   = errors.New("invalid product name")
	ErrInvalidPrice         = errors.New("invalid price")
	ErrInvalidStock         = errors.New("invalid stock")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidCategoryID    = errors.New("invalid category ID")
	ErrDatabaseConnection   = errors.New("database connection error")
	ErrCacheConnection      = errors.New("cache connection error")
)
