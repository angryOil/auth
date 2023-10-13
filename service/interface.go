package service

import (
	"auth/domain"
	"context"
)

type IUserService interface {
	CreateUser(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (string, error)
	Verify(ctx context.Context, verify string) (bool, error)
}
