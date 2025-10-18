package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("email already in use")

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, name, email, password string) error {
	// existingUser, _ := s.repo.FindByEmail(ctx, email)
	// if existingUser != nil {
	// 	return ErrDuplicateEmail
	// }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	user := &User{
		UserID:       uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return fmt.Errorf("could not save the user: %w", err)
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) error {

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("unable to login user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("unable to login user: %w", err)
	}
	return nil

}
