package rankit

import (
	"context"
)

type CreateUserParam struct {
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
}

type AuthenticateUserParam struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type UserService interface {
	CreateUser(ctx context.Context, p CreateUserParam) (*User, error)
	Authenticate(ctx context.Context, p AuthenticateUserParam) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
}
