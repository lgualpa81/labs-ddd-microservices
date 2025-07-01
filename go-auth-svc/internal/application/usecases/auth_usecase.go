package usecases

import (
	"context"
	"fmt"
	"time"

	"poc-auth-svc/internal/application/dtos"
	"poc-auth-svc/internal/domain/services"
	"poc-auth-svc/internal/domain/valueobjects"

	"github.com/golang-jwt/jwt"
)

type AuthUseCase interface {
	Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AuthResponse, error)
	Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.AuthResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*dtos.ValidateResponse, error)
}

type authUseCase struct {
	authService services.AuthService
	jwt         JwtWrapper
}

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

func NewAuthUseCase(authService services.AuthService, config JwtWrapper) AuthUseCase {
	return &authUseCase{
		authService: authService,
		jwt:         config,
	}
}

// Login implements AuthUseCase.
func (uc *authUseCase) Login(ctx context.Context, req *dtos.LoginRequest) (*dtos.AuthResponse, error) {
	user, err := uc.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := uc.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &dtos.AuthResponse{
		Token: token,
		User: &dtos.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

// Register implements AuthUseCase.
func (uc *authUseCase) Register(ctx context.Context, req *dtos.RegisterRequest) (*dtos.AuthResponse, error) {
	user, err := uc.authService.Register(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		return nil, err
	}
	token, err := uc.generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return &dtos.AuthResponse{
		Token: token,
		User: &dtos.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

// ValidateToken implements AuthUseCase.
func (uc *authUseCase) ValidateToken(ctx context.Context, tokenString string) (*dtos.ValidateResponse, error) {
	fmt.Println("usecase secret: " + uc.jwt.SecretKey)
	fmt.Println(uc.jwt.Issuer)
	fmt.Println(uc.jwt.ExpirationHours)
	token, err := jwt.ParseWithClaims(tokenString, &valueobjects.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.jwt.SecretKey), nil
	})
	if err != nil {
		return &dtos.ValidateResponse{Valid: false}, err
	}
	if claims, ok := token.Claims.(*valueobjects.JWTClaims); ok && token.Valid {
		// Opcionalmente verificar si el usuario aún existe y está activo
		user, err := uc.authService.GetUserByID(ctx, claims.UserID)
		if err != nil || !user.IsActive {
			return &dtos.ValidateResponse{Valid: false}, nil
		}

		return &dtos.ValidateResponse{
			Valid: true,
			User: &dtos.UserResponse{
				ID:       user.ID,
				Email:    user.Email,
				Role:     user.Role,
				IsActive: user.IsActive,
			},
			Claims: map[string]interface{}{
				"user_id": claims.UserID,
				"email":   claims.Email,
				"role":    claims.Role,
				//"exp":     claims.ExpiresAt.Unix(),
			},
		}, nil
	}
	return &dtos.ValidateResponse{Valid: false}, nil
}

func (uc *authUseCase) generateToken(userID, email, role string) (signedToken string, err error) {
	claims := &valueobjects.JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(uc.jwt.ExpirationHours)).Unix(),
			Issuer:    uc.jwt.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(uc.jwt.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
