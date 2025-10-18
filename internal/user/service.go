package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo           UserRepository
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

func NewUserService(repo UserRepository, jwtSecret []byte, accessTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		jwtSecret:      jwtSecret,
		accessTokenTTL: accessTokenTTL,
	}
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

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("unable to login user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("unable to login user: %w", err)
	}
	return token, nil

}

func (s *UserService) generateAccessToken(user *User) (string, error) {
	expirationTime := time.Now().Add(s.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.UserID.String(),
		"username": user.Name,
		"email":    user.Email,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}
