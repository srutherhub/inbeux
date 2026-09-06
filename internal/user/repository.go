package user

import (
	"context"
	"errors"
	"fmt"
	"inbeux/internal/apperrors"
	userdb "inbeux/internal/user/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	queries *userdb.Queries
	conn    *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn, queries: userdb.New(conn)}
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (userdb.User, error) {
	if email == "" {
		return userdb.User{}, apperrors.ErrInvalidEmail
	}

	user, err := ur.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userdb.User{}, apperrors.ErrUserNotFound
		}
		return userdb.User{}, fmt.Errorf("get user by email query failed %w", err)
	}

	return user, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, email string) (userdb.User, error) {
	if email == "" {
		return userdb.User{}, apperrors.ErrInvalidEmail
	}

	user, err := ur.queries.CreateUser(ctx, email)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return userdb.User{}, apperrors.ErrUserAlreadyExists
			}
		}
		return userdb.User{}, fmt.Errorf("failed to create user %w", err)
	}

	return user, nil
}
