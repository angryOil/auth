package service

import (
	"auth/domain"
	"context"
)

type IUserService[T any] interface {
	CreateUser(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (T, error)
	Verify(ctx context.Context, verify T) error
}
