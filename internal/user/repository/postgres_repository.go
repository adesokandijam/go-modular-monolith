package repository

import (
	"context"
	"database/sql"
	users "dijam-ecommerce/internal/user"
	"fmt"
)

type PostgresRepository struct {
	DB *sql.DB
}

var _ users.UserRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (p *PostgresRepository) Save(ctx context.Context, user *users.User) error {
	query := `INSERT into users (id, name, email, password_hash)
			VALUES($1,$2,$3,$4);
			`
	_, err := p.DB.ExecContext(ctx, query, user.UserID, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return users.ErrDuplicateEmail
		default:
			return fmt.Errorf("error saving user in postgres db: %w", err)
		}

	}
	return nil
}

func (p *PostgresRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	var user users.User
	query := `SELECT email, name, password_hash FROM users
			WHERE email = $1`
	err := p.DB.QueryRowContext(ctx, query, email).Scan(&user.Email, &user.Name, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

func (p *PostgresRepository) Login(ctx context.Context, email, passwordHash string) (*users.User, error) {

	return nil, nil
}
