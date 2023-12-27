package service

import (
	"context"
	"fmt"

	"rankit/errors"
	"rankit/postgres"
	"rankit/postgres/sqlgen"
	"rankit/rankit"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLoginDetails = errors.E(errors.Invalid, "invalid login details")
	ErrEmailAlreadyExists  = errors.E(errors.Invalid, "email already exists")
	ErrUserNotFound        = errors.E(errors.NotFound, "user not found")
)

const BCRYPT_COST = bcrypt.DefaultCost + 4

type UserService struct {
	querier sqlgen.Querier
}

var _ rankit.UserService = (*UserService)(nil)

func NewUserService(querier sqlgen.Querier) *UserService {
	return &UserService{querier: querier}
}

func (s *UserService) CreateUser(ctx context.Context, p rankit.CreateUserParm) (*rankit.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), BCRYPT_COST)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.querier.CreateUser(ctx, sqlgen.CreateUserParams{
		ID:           generateUserID(),
		Email:        p.Email,
		DisplayName:  p.DisplayName,
		PasswordHash: string(hash),
	})
	if _, ok := postgres.IsUniqueViolation(err); ok {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toRankitUser(user), nil
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*rankit.User, error) {
	user, err := s.querier.GetUserByEmail(ctx, email)
	if err == pgx.ErrNoRows {
		return nil, ErrInvalidLoginDetails
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	isValid := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
	if !isValid {
		return nil, ErrInvalidLoginDetails
	}

	return toRankitUser(user), nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*rankit.User, error) {
	user, err := s.querier.GetUserByID(ctx, id)
	if postgres.IsNotFound(err) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return toRankitUser(user), nil
}

func generateUserID() string {
	return "usr_" + ksuid.New().String()
}

func toRankitUser(u *sqlgen.User) *rankit.User {
	return &rankit.User{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
}
