package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` //no se expone en json
	Role      string    `json:"role" bson:"role"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"" bson:"created_at"`
	UpdatedAt time.Time `json:"" bson:"updated_at"`
}

func NewUser(email, password, role string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	if role == "" {
		role = "user"
	}

	return &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

func (u *User) UpdatePassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return nil
}
