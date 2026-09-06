package user

import (
	"context"

	userdb "inbeux/internal/user/db"
)

type UserService struct {
	query *UserRepository
}

func NewService(repository *UserRepository) *UserService {
	return &UserService{query: repository}
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (userdb.User, error) {
	user, err := us.query.GetUserByEmail(ctx, email)

	if err != nil {
		return userdb.User{}, err
	}

	return user, nil

}

func (us *UserService) CreateUser(ctx context.Context, email string) (userdb.User, error) {
	user, err := us.query.CreateUser(ctx, email)

	if err != nil {
		return userdb.User{}, err
	}

	return user, nil
}
