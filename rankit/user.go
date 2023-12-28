package rankit

import (
	"context"
)

type CreateUserParam struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type AuthenticateUserParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
