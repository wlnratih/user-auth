package repository

import "context"

type RepositoryInterface interface {
	InsertUser(ctx context.Context, input InsertUser) (int, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	InsertUserLoginHistory(ctx context.Context, input InsertUserLoginHistory) error
	GetUserByID(ctx context.Context, id int) (User, error)
	UpdateUser(ctx context.Context, input UpdateUser) (User, error)
}
