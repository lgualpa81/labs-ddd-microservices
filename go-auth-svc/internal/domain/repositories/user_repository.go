package repositories

import (
	"context"
	"errors"

	"poc-auth-svc/internal/domain/entities"
)

var (
	// ErrUserNotFound se devuelve cuando no se encuentra un usuario
	ErrUserNotFound = errors.New("user not found")
	// ErrDuplicateEmail se devuelve cuando el email ya existe
	ErrDuplicateEmail = errors.New("email already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByID(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}
