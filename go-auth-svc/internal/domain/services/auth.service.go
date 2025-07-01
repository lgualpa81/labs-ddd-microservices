package services

import (
	"context"
	"errors"

	"poc-auth-svc/internal/domain/entities"
	err_domain "poc-auth-svc/internal/domain/errors"
	"poc-auth-svc/internal/domain/repositories"
)

type AuthService interface {
	Register(ctx context.Context, email, password, role string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (*entities.User, error)
	GetUserByID(ctx context.Context, id string) (*entities.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type authService struct {
	userRepo repositories.UserRepository
	hasher   PasswordHasher
}

func NewAuthService(userRepo repositories.UserRepository, hasher PasswordHasher) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *authService) Register(ctx context.Context, email, password, role string) (*entities.User, error) {

	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New(err_domain.GetMessage(err_domain.UserAlreadyExists)) //repositories.ErrDuplicateEmail
	}
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := entities.NewUser(email, hashedPassword, role)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New(err_domain.GetMessage(err_domain.InvalidCredentials))
	}
	if !user.IsActive {
		return nil, errors.New(err_domain.GetMessage(err_domain.UserInactive))
	}
	if ok := s.hasher.Compare(user.Password, password); !ok {
		return nil, errors.New(err_domain.GetMessage(err_domain.InvalidCredentials))
	}
	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
