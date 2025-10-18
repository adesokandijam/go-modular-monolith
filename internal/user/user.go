package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	// Login(ctx context.Context, email, passwordHash string) (*User, error)
	// Findby(ctx context.Context, emai, passwordHash string) (*User, error)
}
